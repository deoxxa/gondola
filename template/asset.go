package template

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"gnd.la/log"
	"gnd.la/template/assets"
	"hash/fnv"
	"io"
	"path"
)

func executeAsset(t *Template, m *assets.Manager, asset *assets.Asset) (string, error) {
	name := asset.TemplateName()
	log.Debugf("executing asset template %s", name)
	tmpl := New(m.Loader(), nil)
	tmpl.Funcs(t.funcMap)
	if err := tmpl.Parse(name); err != nil {
		return "", err
	}
	if err := tmpl.Compile(); err != nil {
		return "", err
	}
	vars := t.vars
	if asset.Namespace != "" && asset.Namespace != t.Namespace {
		vars = vars.Unpack(asset.Namespace)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", err
	}
	ext := path.Ext(name)
	nonExt := name[:len(name)-len(ext)]
	f, _, err := m.Load(name)
	if err != nil {
		return "", err
	}
	h := fnv.New32a()
	// Writing to a hash.Hash never fails
	// Include the initial and the final
	// asset in the hash.
	io.Copy(h, f)
	f.Close()
	h.Write(buf.Bytes())
	hash := hex.EncodeToString(h.Sum(nil))
	out := fmt.Sprintf("%s.gen.%s%s", nonExt, hash, ext)
	w, err := m.Create(out, true)
	if err != nil {
		return "", err
	}
	if _, err := w.Write(buf.Bytes()); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	return out, nil
}

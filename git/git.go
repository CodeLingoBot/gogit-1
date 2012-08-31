package git

// #cgo pkg-config: libgit2
// #include <git2.h>
// #include <git2/errors.h>
import "C"

import (
	"errors"
	"unsafe"
)

type Repository struct {
	repository *C.struct_git_repository
}

// NewRepository inits a new repository.
func NewRepository(path string, bare bool) (*Repository, error) {
	var cbare C.unsigned = 0
	if bare {
		cbare = 1
	}
	repo := &Repository{}
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	if C.git_repository_init(&repo.repository, cpath, cbare) != C.GIT_OK {
		return nil, LastErr()
	}
	return repo, nil
}

// GetRepository opens a repository by its path.
func GetRepository(path string) (*Repository, error) {
	repo := &Repository{}
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	if C.git_repository_open(&repo.repository, cpath) != C.GIT_OK {
		return nil, LastErr()
	}
	return repo, nil
}

func (r *Repository) Config() (*Config, error) {
	conf := &Config{}
	if C.git_repository_config(&conf.config, r.repository) != C.GIT_OK {
		return nil, LastErr()
	}
	return conf, nil
}

func (r *Repository) Free() {
	C.git_repository_free(r.repository)
}

type Config struct {
	config *C.struct_git_config
}

func (c *Config) Free() {
	C.git_config_free(c.config)
}

func (c *Config) GetBool(name string) (bool, error) {
	var v C.int
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	if C.git_config_get_bool(&v, c.config, cname) != C.GIT_OK {
		return false, LastErr()
	}
	return v == 1, nil
}

func (c *Config) SetBool(name string, value bool) error {
	var v C.int = 0
	if value {
		v = 1
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	if C.git_config_set_bool(c.config, cname, v) != C.GIT_OK {
		return LastErr()
	}
	return nil
}

func (c *Config) GetString(name string) (string, error) {
	var v *C.char
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	if C.git_config_get_string(&v, c.config, cname) != C.GIT_OK {
		return "", LastErr()
	}
	return C.GoString(v), nil
}

func (c *Config) SetString(name, value string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	if C.git_config_set_string(c.config, cname, cvalue) != C.GIT_OK {
		return LastErr()
	}
	return nil
}

func (c *Config) GetInt64(name string) (int64, error) {
	var v C.int64_t
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	if C.git_config_get_int64(&v, c.config, cname) != C.GIT_OK {
		return 0, LastErr()
	}
	return int64(v), nil
}

func (c *Config) SetInt64(name string, value int64) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	if C.git_config_set_int64(c.config, cname, C.int64_t(value)) != C.GIT_OK {
		return LastErr()
	}
	return nil
}

func LastErr() error {
	err := C.giterr_last()
	return errors.New(C.GoString(err.message))
}

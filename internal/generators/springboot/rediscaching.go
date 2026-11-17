package springboot

import "strings"

type RedisCachingConfig struct {
	pg projectGenerator
}

func NewRedisCachingConfig(pg projectGenerator) *RedisCachingConfig {
	return &RedisCachingConfig{pg: pg}
}

func (r RedisCachingConfig) generate(pc ProjectConfig) error {
	if !pc.RedisCachingSupport {
		return nil
	}
	if err := r.createSrcMainJava(pc); err != nil {
		return err
	}
	return nil
}

func (r RedisCachingConfig) createSrcMainJava(pc ProjectConfig) error {
	basePackagePath := strings.ReplaceAll(pc.BasePackage, ".", "/")

	templateMap := map[string]string{
		"config/CacheConfig.java.tmpl": "config/CacheConfig.java",
	}

	for tmpl, filePath := range templateMap {
		err := r.pg.executeTemplate(pc, srcMainJavaPath+tmpl, srcMainJavaPath+basePackagePath+"/"+filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

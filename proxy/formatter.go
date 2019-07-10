/**
 * @Author: huangw1
 * @Date: 2019/7/11 15:42
 */

package proxy

import (
	"strings"
)

type EntityFormatter interface {
	Format(entity Response) Response
}

type propertyFilter func(entity *Response)

type entityFormatter struct {
	Target         string
	Prefix         string
	PropertyFilter propertyFilter
	Mapping        map[string]string
}

func (e entityFormatter) Format(entity Response) Response {
	if e.Target != "" {
		extractTarget(e.Target, &entity)
	}
	if len(entity.Data) > 0 {
		e.PropertyFilter(&entity)
	}
	if len(entity.Data) > 0 {
		for key, newKey := range e.Mapping {
			if v, ok := entity.Data[key]; ok {
				entity.Data[newKey] = v
				delete(entity.Data, key)
			}
		}
	}
	if e.Prefix != "" {
		entity.Data = map[string]interface{}{e.Prefix: entity.Data}
	}
	return entity
}

func extractTarget(target string, entity *Response) {
	if temp, ok := entity.Data[target]; ok {
		entity.Data, ok = temp.(map[string]interface{})
		if !ok {
			entity.Data = map[string]interface{}{}
		}
	} else {
		entity.Data = map[string]interface{}{}
	}
}

func NewEntityFormatter(target string, whitelist, blacklist []string, group string, mappings map[string]string) EntityFormatter {
	var propertyFilter propertyFilter
	if len(whitelist) > 0 {
		propertyFilter = newWhitelistingFilter(whitelist)
	} else {
		propertyFilter = newBlacklistFilter(blacklist)
	}
	sanitizedMappings := make(map[string]string, len(mappings))
	for k, v := range mappings {
		sanitizedMappings[k] = strings.Split(v, ".")[0]
	}
	return entityFormatter{
		Target:         target,
		Prefix:         group,
		PropertyFilter: propertyFilter,
		Mapping:        sanitizedMappings,
	}
}

func newWhitelistingFilter(whitelist []string) propertyFilter {
	wl := make(map[string]map[string]interface{}, len(whitelist))
	for _, k := range whitelist {
		keys := strings.Split(k, ".")
		temp := make(map[string]interface{}, len(keys) - 1)
		if len(keys) > 1 {
			if _, ok := wl[keys[0]]; ok {
				for _, key := range keys[1:] {
					wl[keys[0]][key] = nil
				}
			} else {
				for _, key := range keys[1:] {
					temp[key] = nil
				}
				wl[keys[0]] = temp
			}
		} else {
			wl[keys[0]] = temp
		}
	}
	return func(entity *Response) {
		accumulator := make(map[string]interface{}, len(whitelist))
		for k, v := range entity.Data {
			if sub, ok := wl[k]; ok {
				if len(sub) > 0 {
					if temp := whitelistFilterSub(v, sub); len(temp) > 0 {
						accumulator[k] = temp
					}
				} else {
					accumulator[k] = v
				}
			}
		}
	}
}

func whitelistFilterSub(v interface{}, whitelist map[string]interface{}) map[string]interface{} {
	entity, ok := v.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
	}
	temp := make(map[string]interface{}, len(whitelist))
	for k, v := range entity {
		if _, ok := whitelist[k]; ok {
			temp[k] = v
		}
	}
	return temp
}

func newBlacklistFilter(blacklist []string) propertyFilter {
	bl := make(map[string][]string, len(blacklist))
	for _, k := range blacklist {
		keys := strings.Split(k, ".")
		if len(keys) > 1 {
			if sub, ok := bl[keys[0]]; ok {
				bl[keys[0]] = append(sub, keys[1])
			} else {
				bl[keys[0]] = []string{keys[1]}
			}
		} else {
			bl[keys[0]] = []string{}
		}
	}
	return func(entity *Response) {
		for k, v := range bl {
			if len(v) == 0 {
				delete(entity.Data, k)
			} else {
				if temp := blacklistFilterSub(entity.Data, v); len(temp) > 0 {
					entity.Data[k] = temp
				}
			}
		}
	}
}

func blacklistFilterSub(v interface{}, sub []string) map[string]interface{} {
	entity, ok := v.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
	}
	for _, key := range sub {
		delete(entity, key)
	}
	return entity
}

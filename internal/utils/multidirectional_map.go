package utils

import (
	"errors"
	"strconv"
)

type Scene int

var OriginalScene = Scene(-1)

type MultidirectionalMap struct {
	mapping *map[Scene][]any
}

// Map 映射，把源场景的值映射到目标场景的值
func (m MultidirectionalMap) Map(source any, sourceScene Scene, targetScene Scene) (any, error) {
	sourceValues, ok := (*m.mapping)[sourceScene]
	if !ok {
		return nil, errors.New("source scene: " + strconv.Itoa(int(sourceScene)) + " not define")
	}
	targetValues, ok := (*m.mapping)[targetScene]
	if !ok {
		return nil, errors.New("target scene: " + strconv.Itoa(int(targetScene)) + "not define")
	}
	for i, value := range sourceValues {
		if source == value {
			return targetValues[i], nil
		}
	}
	return nil, errors.New("source value not found in " + strconv.Itoa(int(sourceScene)))
}

// FromOriginalScene 从OriginalScene映射到目标场景
func (m MultidirectionalMap) FromOriginalScene(source any, targetScene Scene) (any, error) {
	return m.Map(source, OriginalScene, targetScene)
}

// ToOriginalScene 从源场景映射到OriginalScene
func (m MultidirectionalMap) ToOriginalScene(source any, sourceScene Scene) (any, error) {
	return m.Map(source, sourceScene, OriginalScene)
}

// NewMultidirectionalMap 新建映射
func NewMultidirectionalMap(mapping map[Scene][]any) (*MultidirectionalMap, error) {
	if mapping == nil {
		return nil, errors.New("MultidirectionalMap mapping cannot be nil")
	}
	if _, ok := mapping[OriginalScene]; !ok {
		return nil, errors.New("MultidirectionalMap mapping must contain OriginalScene")
	}
	count := -1
	for scene, values := range mapping {
		if count == -1 {
			count = len(values)
		} else {
			if len(values) != count {
				return nil, errors.New("scene [" + strconv.Itoa(int(scene)) + "] has: " + strconv.Itoa(len(values)) + " elements, but target has: " + strconv.Itoa(count))
			}
		}
		ok, i, j := hasDuplicates(values)
		if ok {
			return nil, errors.New("scene [" + strconv.Itoa(int(scene)) + "] has duplicate elements, index of: " + strconv.Itoa(i) + ", " + strconv.Itoa(j))
		}
	}
	return &MultidirectionalMap{mapping: &mapping}, nil
}

// Must 新建映射
func Must(originalScenes []any, mapping map[Scene][]any) *MultidirectionalMap {
	mapping[OriginalScene] = originalScenes
	multidirectionalMap, err := NewMultidirectionalMap(mapping)
	if err != nil {
		panic(err)
	}
	return multidirectionalMap
}

// hasDuplicates 校验是否存在重复元素
func hasDuplicates(slice []any) (bool, int, int) {
	for i := 0; i < len(slice); i++ {
		for j := i + 1; j < len(slice); j++ {
			if slice[i] == slice[j] {
				return true, i, j
			}
		}
	}
	return false, -1, -1
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Obj struct {
	Faces []Face

	vertices []Vertex3
	textures []Vertex2
	normals  []Vertex3
}

func loadObjFromFile(filename string) (*Obj, error) {
	obj := Obj{}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lineNumber := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")

		switch parts[0] {
		// Vertex4 line
		case "v":
			if err := obj.parseVertexLine(line, lineNumber); err != nil {
				return nil, err
			}

		// Vertex4 normal line
		case "vn":
			if err := obj.parseVertexNormalLine(line, lineNumber); err != nil {
				return nil, err
			}

		// Vertex4 texture line
		case "vt":
			if err := obj.parseVertexTextureLine(line, lineNumber); err != nil {
				return nil, err
			}

		// Face line
		case "f":
			if err := obj.parseFaceLine(line, lineNumber); err != nil {
				return nil, err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Cleanup
	obj.vertices = []Vertex3{}
	obj.normals = []Vertex3{}
	obj.textures = []Vertex2{}

	return &obj, nil
}

func (obj *Obj) parseFaceLine(line string, lineNumber int) error {
	parts := strings.Split(line, " ")

	// Remove empty parts
	for k, v := range parts {
		if v == "" {
			parts = append(parts[:k], parts[k+1:]...)
		}
	}

	if len(parts) < 4 {
		return errors.New(fmt.Sprintf("insufficient points found in face directive on line %d", lineNumber))
	}

	// First vertex
	firstArgs := strings.Split(parts[1], "/")
	vertexId, err := strconv.Atoi(firstArgs[0])
	if err != nil {
		return err
	}
	vertexTextureId, err := strconv.Atoi(firstArgs[1])
	if err != nil {
		return err
	}
	vertexNormalId, err := strconv.Atoi(firstArgs[2])
	if err != nil {
		return err
	}
	firstVertex, err := obj.resolveVertexId(vertexId, lineNumber)
	if err != nil {
		return err
	}
	firstVertexTexture, err := obj.resolveVertexTextureId(vertexTextureId, lineNumber)
	if err != nil {
		return err
	}
	firstVertexNormal, err := obj.resolveVertexNormalId(vertexNormalId, lineNumber)
	if err != nil {
		return err
	}

	// Second vertex
	secondArgs := strings.Split(parts[2], "/")
	vertexId, err = strconv.Atoi(secondArgs[0])
	if err != nil {
		return err
	}
	vertexTextureId, err = strconv.Atoi(secondArgs[1])
	if err != nil {
		return err
	}
	vertexNormalId, err = strconv.Atoi(secondArgs[2])
	if err != nil {
		return err
	}
	secondVertex, err := obj.resolveVertexId(vertexId, lineNumber)
	if err != nil {
		return err
	}
	secondVertexTexture, err := obj.resolveVertexTextureId(vertexTextureId, lineNumber)
	if err != nil {
		return err
	}
	secondVertexNormal, err := obj.resolveVertexNormalId(vertexNormalId, lineNumber)
	if err != nil {
		return err
	}

	// Third vertex
	thirdArgs := strings.Split(parts[3], "/")
	vertexId, err = strconv.Atoi(thirdArgs[0])
	if err != nil {
		return err
	}
	vertexTextureId, err = strconv.Atoi(thirdArgs[1])
	if err != nil {
		return err
	}
	vertexNormalId, err = strconv.Atoi(thirdArgs[2])
	if err != nil {
		return err
	}
	thirdVertex, err := obj.resolveVertexId(vertexId, lineNumber)
	if err != nil {
		return err
	}
	thirdVertexTexture, err := obj.resolveVertexTextureId(vertexTextureId, lineNumber)
	if err != nil {
		return err
	}
	thirdVertexNormal, err := obj.resolveVertexNormalId(vertexNormalId, lineNumber)
	if err != nil {
		return err
	}

	obj.Faces = append(obj.Faces, Face{
		Vertices: [3]Vertex3{
			firstVertex,
			secondVertex,
			thirdVertex,
		},
		Textures: [3]Vertex2{
			firstVertexTexture,
			secondVertexTexture,
			thirdVertexTexture,
		},
		Normals: [3]Vertex3{
			firstVertexNormal,
			secondVertexNormal,
			thirdVertexNormal,
		},
	})

	return nil
}

func (obj *Obj) resolveVertexId(id int, lineNumber int) (Vertex3, error) {
	if id > len(obj.vertices) {
		return Vertex3{}, errors.New(fmt.Sprintf("unable to resolve vertex id %d used on line %d", id, lineNumber))
	}
	return obj.vertices[id-1], nil
}

func (obj *Obj) resolveVertexNormalId(id int, lineNumber int) (Vertex3, error) {
	if id > len(obj.normals) {
		return Vertex3{}, errors.New(fmt.Sprintf("unable to resolve vertex normal id %d used on line %d", id, lineNumber))
	}
	return obj.normals[id-1], nil
}

func (obj *Obj) resolveVertexTextureId(id int, lineNumber int) (Vertex2, error) {
	if id > len(obj.textures) {
		return Vertex2{}, errors.New(fmt.Sprintf("unable to resolve vertex texture id %d used on line %d", id, lineNumber))
	}
	return obj.textures[id-1], nil
}

func (obj *Obj) parseVertexLine(line string, lineNumber int) error {
	vertex := Vertex3{}

	parts := strings.Split(line, " ")

	if len(parts) < 4 {
		return errors.New(fmt.Sprintf("insufficient points found in vertex directive on line %d", lineNumber))
	}

	// Remove empty parts
	for k, v := range parts {
		if v == "" {
			parts = append(parts[:k], parts[k+1:]...)
		}
	}

	var err error

	// X
	vertex.X, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid float x coordinate found in vertex directive on line %d", lineNumber))
	}

	// Y
	vertex.Y, err = strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid float y coordinate found in vertex directive on line %d", lineNumber))
	}

	// Z
	vertex.Z, err = strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid float z coordinate found in vertex directive on line %d", lineNumber))
	}

	obj.vertices = append(obj.vertices, vertex)

	return nil
}

func (obj *Obj) parseVertexNormalLine(line string, lineNumber int) error {
	vertexNormal := Vertex3{}

	parts := strings.Split(line, " ")

	if len(parts) < 4 {
		return errors.New(fmt.Sprintf("insufficient points found in vertex normal directive on line %d", lineNumber))
	}

	// Remove empty parts
	for k, v := range parts {
		if v == "" {
			parts = append(parts[:k], parts[k+1:]...)
		}
	}

	var err error

	// X
	vertexNormal.X, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid float x coordinate found in vertex normal directive on line %d", lineNumber))
	}

	// Y
	vertexNormal.Y, err = strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid float y coordinate found in vertex normal directive on line %d", lineNumber))
	}

	// Z
	vertexNormal.Z, err = strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid float z coordinate found in vertex normal directive on line %d", lineNumber))
	}

	obj.normals = append(obj.normals, vertexNormal)

	return nil
}

func (obj *Obj) parseVertexTextureLine(line string, lineNumber int) error {
	vertexTexture := Vertex2{}

	parts := strings.Split(line, " ")

	if len(parts) < 3 {
		return errors.New(fmt.Sprintf("insufficient points found in vertex texture directive on line %d", lineNumber))
	}

	// Remove empty parts
	for k, v := range parts {
		if v == "" {
			parts = append(parts[:k], parts[k+1:]...)
		}
	}

	var err error

	// X
	vertexTexture.X, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid float x coordinate found in vertex texture directive on line %d", lineNumber))
	}

	// Y
	vertexTexture.Y, err = strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return errors.New(fmt.Sprintf("invalid float y coordinate found in vertex texture directive on line %d", lineNumber))
	}

	obj.textures = append(obj.textures, vertexTexture)

	return nil
}

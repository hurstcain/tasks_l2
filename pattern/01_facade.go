/*
Паттерн Фасад используется для того, чтобы предоставлять упрощенный интерфейс
для работы со сложными подсистемами, а также изолировать клиента от компонентов сложной подсистемы.
То есть Фасад создает новый упрощенный интерфейс на основе большой подсистемы, которая в свою очередь
может состоять из нескольких интерфейсов.

Плюсы:
- упрощает код, с которым работает клиент;
- изолирует клиента от компонентов сложной подсистемы.
Минусы:
- может начать разрастаться и в итоге стать божественным объектом. Но чтобы этого не случилось,
можно реализовывать для сложной подсистемы вместо одного большого фасада - несколько простых.

Примеры использования:
Например, есть библиотека для работы с видеофайлами, состоящая из большого количества компонентов.
Чтобы упростить работу с этой библиотекой для клиента, можно создать структуру-фасад, которая для
работы будет использовать компоненты исходной библиотеки.
*/

package pattern

import (
	"log"
	"os"
	"strings"
)

// Сложная подсистема, которая используется для работы с видофайлами и конвертации видеофайлов.
// Для этой подсистемы будет реализован фасад.

// Codec - интерфейс кодека.
type Codec interface {
	CType() string
}

// VideoFile - структура видеофайла.
type VideoFile struct {
	name      string
	codecType string
}

// NewVideoFile - конструктор структуры VideoFile.
func NewVideoFile(name string) VideoFile {
	cType := make([]rune, 0)
	index := strings.LastIndexByte(name, '.') + 1
	cType = append(cType, []rune(name)[index:]...)

	return VideoFile{
		name:      name,
		codecType: string(cType),
	}
}

// Name - возвращает имя видеофайла.
func (vf VideoFile) Name() string {
	return vf.name
}

// CodecType - возвращает тип кодека видеофайла.
func (vf VideoFile) CodecType() string {
	return vf.codecType
}

// MPEG4CompressionCodec - структура кодека MPEG4.
type MPEG4CompressionCodec struct {
	ctype string
}

// NewMPEG4CompressionCodec - конструктор структуры MPEG4CompressionCodec.
func NewMPEG4CompressionCodec() MPEG4CompressionCodec {
	return MPEG4CompressionCodec{
		ctype: "mp4",
	}
}

// CType - возвращает тип кодека.
func (c MPEG4CompressionCodec) CType() string {
	return c.ctype
}

// OggCompressionCodec - структура кодека Ogg.
type OggCompressionCodec struct {
	ctype string
}

// NewOggCompressionCodec - конструктор структуры OggCompressionCodec.
func NewOggCompressionCodec() OggCompressionCodec {
	return OggCompressionCodec{
		ctype: "ogg",
	}
}

// CType - возвращает тип кодека.
func (c OggCompressionCodec) CType() string {
	return c.ctype
}

// CodecFactory - фабрика видеокодеков.
type CodecFactory struct{}

// Extract - извлечение кодека видеофайла.
func (cf CodecFactory) Extract(file VideoFile) Codec {
	if file.CodecType() == "mp4" {
		log.Printf("CodecFactory: file %s, extracting mpeg audio...\n", file.Name())
		return NewMPEG4CompressionCodec()
	} else {
		log.Printf("CodecFactory: file %s, extracting ogg audio...\n", file.Name())
		return NewOggCompressionCodec()
	}
}

// BitrateReader - Bitrate-конвертер.
type BitrateReader struct{}

// Read - чтение видеофайла.
func (br BitrateReader) Read(file VideoFile, codec Codec) VideoFile {
	log.Printf("BitrateReader: reading file %s, codec: %s...\n", file.Name(), codec.CType())
	return file
}

// Convert - конвертация видеофайла.
func (br BitrateReader) Convert(buffer VideoFile, codec Codec) VideoFile {
	log.Printf("BitrateReader: convert file %s, codec: %s...\n", buffer.Name(), codec.CType())
	var fname []rune
	index := strings.LastIndexByte(buffer.name, '.')
	fname = append(fname, []rune(buffer.name)[:index]...)
	fname = append(fname, '1', '.')
	fname = append(fname, []rune(codec.CType())...)
	newVideoFile := VideoFile{
		name:      string(fname),
		codecType: codec.CType(),
	}
	log.Printf("BitrateReader: new file: %s\n", newVideoFile.Name())

	return newVideoFile
}

// AudioMixer - структура микширования аудио.
type AudioMixer struct{}

// Fix - микширования аудио.
func (am AudioMixer) Fix(result VideoFile) *os.File {
	log.Printf("AudioMixer: fixing audio...")
	file, _ := os.Create(result.name)
	return file
}

// Фасад для библиотеки для работы с видеофайлами.
// Упрощает для клиента конвертацию видео и скрывает элементы исходной подсистемы.
// Фасад состоит только из одной функции, которая конвертирует видеофайлы.

// VideoConversionFacade - структура конвертации видеофайлов.
type VideoConversionFacade struct {
	// Компоненты подсистемы, которые будут использоваться в фасаде.
	codecFactory  CodecFactory
	bitrateReader BitrateReader
	audioMixer    AudioMixer
}

// ConvertVideo - функция конвертации видеофайла.
func (vcf VideoConversionFacade) ConvertVideo(fileName string, codec string) *os.File {
	log.Println("VideoConversionFacade: conversion started...")
	file := NewVideoFile(fileName)
	fileCodec := vcf.codecFactory.Extract(file)
	var destinationCodec Codec
	if codec == "mp4" {
		destinationCodec = NewMPEG4CompressionCodec()
	} else {
		destinationCodec = NewOggCompressionCodec()
	}
	buffer := vcf.bitrateReader.Read(file, fileCodec)
	intermediateResult := vcf.bitrateReader.Convert(buffer, destinationCodec)
	result := vcf.audioMixer.Fix(intermediateResult)
	log.Println("VideoConversionFacade: conversion completed.")

	return result
}

// FacadeClient - код клиента
func FacadeClient() {
	converter := VideoConversionFacade{}
	file := converter.ConvertVideo("video.mp4", "ogg")
	log.Println(file.Name())
	/*
		output:
		2022/02/28 12:50:12 VideoConversionFacade: conversion started...
		2022/02/28 12:50:12 CodecFactory: file video.mp4, extracting mpeg audio...
		2022/02/28 12:50:12 BitrateReader: reading file video.mp4, codec: mp4...
		2022/02/28 12:50:12 BitrateReader: convert file video.mp4, codec: ogg...
		2022/02/28 12:50:12 BitrateReader: new file: video1.ogg
		2022/02/28 12:50:12 AudioMixer: fixing audio...
		2022/02/28 12:50:12 VideoConversionFacade: conversion completed.
		2022/02/28 12:50:12 video1.ogg
	*/
}

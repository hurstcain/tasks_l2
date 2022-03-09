/*
Фабричный метод - это порождающий паттерн проектирования, применяется для создания объектов
с определенным интерфейсом, а различные реализации этого интерфейса предоставляются потомкам.
Другими словами, есть базовый абстрактный класс фабрики, который говорит,
что каждая фабрика, наследующая этот класс, должна реализовать такой-то метод для создания своих продуктов.

Плюсы:
- так как создание объекта происходит только в одном классе, то в будущем изменения в коде
нужно будет вносить только в одном конкретном месте;
- упрощается добавление новых продуктов в программу.
Минусы:
- если продуктов очень много, то появляется много классов, которые похожи между собой.
*/

package pattern

import (
	"errors"
	"fmt"
	"os"
)

// В данном примере есть файлы с расширениями .txt, .log и .doc.
// В зависимости от расширения, файлы сохраняются в конкретную папку: /file/txt/, /file/log/ и file/doc/.
// Продукты - это структуры файлов с различными тремя расширениями.
// У каждого из этих трех структур различные значения полей extension (расширение) и path (путь к файлу).
// Эти структуры встраивают основную структуру File, которая содержит метод Save,
// сохраняющий файлы в конкретной директории.
// Фабричный метод в данном случае будет принимать на вход расширение файла и в зависимости от него
// возвращать экземпляр одной из трех структур, описывающих файл.

// FileInterface - интерфейс продукта, который возвращает фабрика.
// В данном случае это некоторый файл.
type FileInterface interface {
	GetName() string
	SetName(string)
	GetFilePath() string
	GetFileExtension() string
	Save()
}

// FactoryInterface - интерфейс фабрики для производства продуктов.
// В данном случае фабрика создает некоторый файл.
type FactoryInterface interface {
	CreateFile(string, string) (FileInterface, error)
}

// File - структура конкретного продукта.
// Реализует интерфейс FileInterface.
// Является базовой структурой файла, в которой описаны все методы.
// Структура встраивается в другие структуры, которые описывают файлы с конкретным расширением.
// name - название файла.
// path - путь, по которому файл сохраняется.
// extension - расширение файла.
type File struct {
	name      string
	path      string
	extension string
}

// GetName - возвращает название файла.
func (f File) GetName() string {
	return f.name
}

// SetName - изменяет название файла.
func (f *File) SetName(s string) {
	f.name = s
}

// GetFilePath - возвращает путь к файлу.
func (f File) GetFilePath() string {
	return f.path
}

// GetFileExtension - возвращает расширение файла.
func (f File) GetFileExtension() string {
	return f.extension
}

// Save - сохраняет файл в некоторую директорию, путь к которой зависит от расширения файла.
func (f File) Save() {
	file, _ := os.Create(f.path + f.name + f.extension)
	file.WriteString(f.path + f.name + f.extension)
	file.Close()
	fmt.Printf("%s saved in %s\n", f.name+f.extension, f.path)
}

// TxtFile - конкретный продукт, описывает файл с расширением .txt.
// В данную структуру встраивается структура File.
type TxtFile struct {
	File
}

// NewTxtFile - конструктор структуры TxtFile.
func NewTxtFile(name string) *TxtFile {
	return &TxtFile{
		File{
			name:      name,
			path:      "file/txt/",
			extension: ".txt",
		},
	}
}

// LogFile - конкретный продукт, описывает файл с расширением .log.
// В данную структуру встраивается структура File.
type LogFile struct {
	File
}

// NewLogFile - конструктор структуры LogFile.
func NewLogFile(name string) *LogFile {
	return &LogFile{
		File{
			name:      name,
			path:      "file/log/",
			extension: ".log",
		},
	}
}

// DocFile - конкретный продукт, описывает файл с расширением .doc.
// В данную структуру встраивается структура File.
type DocFile struct {
	File
}

// NewDocFile - конструктор структуры DocFile.
func NewDocFile(name string) *DocFile {
	return &DocFile{
		File{
			name:      name,
			path:      "file/doc/",
			extension: ".doc",
		},
	}
}

// Factory - фабрика, структура, которая производит продукты.
type Factory struct{}

// CreateFile - фабричный метод.
// В зависимости от расширения файла возвращает экземпляр структуры, описывающей файл с конкретным расширением.
func (f Factory) CreateFile(name, extension string) (FileInterface, error) {
	switch extension {
	case "txt":
		return NewTxtFile(name), nil
	case "log":
		return NewLogFile(name), nil
	case "doc":
		return NewDocFile(name), nil
	default:
		return nil, errors.New("неверное расширение файла")
	}
}

// FactoryMethodClient - код клиента.
func FactoryMethodClient() {
	// Экземпляр фабрики.
	factory := Factory{}

	// Создаем три файла с разными расширениями.
	file1, _ := factory.CreateFile("file", "txt")
	file2, _ := factory.CreateFile("file", "log")
	file3, _ := factory.CreateFile("file", "doc")

	// Вывод информации о созданных файлах.
	fmt.Println("file 1 info:")
	fmt.Printf("name: %s\npath: %s\n", file1.GetName()+file1.GetFileExtension(), file1.GetFilePath())
	fmt.Println("file 2 info:")
	fmt.Printf("name: %s\npath: %s\n", file2.GetName()+file2.GetFileExtension(), file2.GetFilePath())
	fmt.Println("file 3 info:")
	fmt.Printf("name: %s\npath: %s\n", file3.GetName()+file3.GetFileExtension(), file3.GetFilePath())
	fmt.Println("Saving files...")
	// Сохраняем файлы.
	file1.Save()
	file2.Save()
	file3.Save()
	// output:
	// file 1 info:
	// name: file.txt
	// path: file/txt/
	// file 2 info:
	// name: file.log
	// path: file/log/
	// file 3 info:
	// name: file.doc
	// path: file/doc/
	// Saving files...
	// file.txt saved in file/txt/
	// file.log saved in file/log/
	// file.doc saved in file/doc/
}

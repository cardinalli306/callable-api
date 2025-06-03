package storage

import (
	"context"
	"fmt"
	"io"
	
	"strings"
	"time"

	"cloud.google.com/go/storage"
)

// CloudStorage representa uma interface com o Cloud Storage do GCP
type CloudStorage struct {
	bucketName string
	// Mapas para simulação
	mockFiles  map[string][]byte
}

// NewCloudStorage cria uma nova instância de CloudStorage
func NewCloudStorage(bucketName string) *CloudStorage {
	return &CloudStorage{
		bucketName: bucketName,
		mockFiles:  make(map[string][]byte),
	}
}

// GetClient retorna um cliente simulado de Cloud Storage
func (cs *CloudStorage) GetClient(ctx context.Context) (*storage.Client, error) {
	// Simular criação de cliente bem-sucedida
	fmt.Printf("[MOCK] Cloud Storage client criado para bucket: %s\n", cs.bucketName)
	return &storage.Client{}, nil
}

// UploadFile simula o upload de um arquivo para o Cloud Storage
func (cs *CloudStorage) UploadFile(ctx context.Context, objectName string, file io.Reader) error {
	// Ler conteúdo do arquivo
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	
	// Armazenar no mapa de simulação
	cs.mockFiles[objectName] = data
	fmt.Printf("[MOCK] Arquivo simulado upload: %s (tamanho: %d bytes)\n", objectName, len(data))
	return nil
}

// DownloadFile simula o download de um arquivo do Cloud Storage
func (cs *CloudStorage) DownloadFile(ctx context.Context, objectName string) ([]byte, error) {
	// Verificar se o arquivo existe no mapa de simulação
	if data, exists := cs.mockFiles[objectName]; exists {
		fmt.Printf("[MOCK] Arquivo simulado download: %s\n", objectName)
		return data, nil
	}
	
	// Se não existir, criar dados simulados
	mockData := []byte(fmt.Sprintf("Conteúdo simulado para %s criado em %s", 
		objectName, time.Now().Format(time.RFC3339)))
	cs.mockFiles[objectName] = mockData
	fmt.Printf("[MOCK] Arquivo simulado criado on-demand: %s\n", objectName)
	return mockData, nil
}

// DeleteFile simula a exclusão de um arquivo do Cloud Storage
func (cs *CloudStorage) DeleteFile(ctx context.Context, objectName string) error {
	// Verificar se o arquivo existe no mapa de simulação
	if _, exists := cs.mockFiles[objectName]; exists {
		delete(cs.mockFiles, objectName)
		fmt.Printf("[MOCK] Arquivo simulado excluído: %s\n", objectName)
		return nil
	}
	
	fmt.Printf("[MOCK] Tentativa de exclusão de arquivo inexistente: %s\n", objectName)
	return nil // Não retornamos erro para simular sucesso
}

// ListFiles simula a listagem de arquivos em um diretório do Cloud Storage
func (cs *CloudStorage) ListFiles(ctx context.Context, prefix string) ([]string, error) {
	var files []string
	
	// Iterar sobre os arquivos simulados
	for key := range cs.mockFiles {
		if strings.HasPrefix(key, prefix) {
			files = append(files, key)
		}
	}
	
	// Se não houver arquivos com esse prefixo, criar alguns para teste
	if len(files) == 0 {
		mockPaths := []string{
			fmt.Sprintf("%sfile1.txt", prefix),
			fmt.Sprintf("%sfile2.pdf", prefix),
			fmt.Sprintf("%ssubdir/file3.json", prefix),
		}
		
		for _, path := range mockPaths {
			cs.mockFiles[path] = []byte(fmt.Sprintf("Conteúdo simulado para %s", path))
			files = append(files, path)
		}
	}
	
	fmt.Printf("[MOCK] Arquivos listados com prefixo '%s': %d arquivos\n", prefix, len(files))
	return files, nil
}

// GetSignedURL simula a geração de uma URL assinada para um objeto
func (cs *CloudStorage) GetSignedURL(ctx context.Context, objectName string, expiration time.Duration) (string, error) {
	mockURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s?mock-signed=true&expires=%d", 
		cs.bucketName, objectName, time.Now().Add(expiration).Unix())
	
	fmt.Printf("[MOCK] URL assinada simulada gerada para: %s\n", objectName)
	return mockURL, nil
}
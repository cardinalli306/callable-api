package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// CloudStorage gerencia operações no Cloud Storage
type CloudStorage struct {
	bucketName string
}

// NewCloudStorage cria uma nova instância do gerenciador de armazenamento
func NewCloudStorage(bucketName string) *CloudStorage {
	return &CloudStorage{
		bucketName: bucketName,
	}
}

// GetClient retorna um novo cliente de armazenamento GCP
func (s *CloudStorage) GetClient(ctx context.Context) (*storage.Client, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar cliente: %v", err)
	}
	return client, nil
}

// GetBucketName retorna o nome do bucket configurado
func (s *CloudStorage) GetBucketName() string {
	return s.bucketName
}

// UploadFile faz upload de um arquivo para o Cloud Storage
func (s *CloudStorage) UploadFile(ctx context.Context, objectName string, data io.Reader) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("falha ao criar cliente: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(s.bucketName)
	obj := bucket.Object(objectName)
	wc := obj.NewWriter(ctx)

	if _, err = io.Copy(wc, data); err != nil {
		return fmt.Errorf("erro na cópia dos dados: %v", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("erro ao fechar writer: %v", err)
	}

	return nil
}

// DownloadFile baixa um arquivo do Cloud Storage
func (s *CloudStorage) DownloadFile(ctx context.Context, objectName string) ([]byte, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar cliente: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(s.bucketName)
	obj := bucket.Object(objectName)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar reader: %v", err)
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler dados: %v", err)
	}

	return data, nil
}

// ListFiles lista arquivos em um diretório do bucket
func (s *CloudStorage) ListFiles(ctx context.Context, prefix string) ([]string, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar cliente: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(s.bucketName)
	var objects []string

	it := bucket.Objects(ctx, &storage.Query{Prefix: prefix})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("erro ao listar objetos: %v", err)
		}
		objects = append(objects, attrs.Name)
	}

	return objects, nil
}

// DeleteFile deleta um arquivo do Cloud Storage
func (s *CloudStorage) DeleteFile(ctx context.Context, objectName string) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("falha ao criar cliente: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(s.bucketName)
	obj := bucket.Object(objectName)

	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("erro ao deletar objeto: %v", err)
	}

	return nil
}

// GetSignedURL gera uma URL assinada para acesso temporário a um objeto
func (s *CloudStorage) GetSignedURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(expires),
		// Removi os campos Bucket e Object que estavam causando o erro
		// Nota: Dependendo do seu ambiente, você pode precisar configurar as credenciais:
		// GoogleAccessID: "seu-service-account-email@projeto.iam.gserviceaccount.com",
		// PrivateKey:     []byte("sua-chave-privada"),
	}

	url, err := storage.SignedURL(s.bucketName, objectName, opts)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar URL assinada: %v", err)
	}

	return url, nil
}
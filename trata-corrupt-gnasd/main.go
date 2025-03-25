package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	diretorio := "/var/opt/mediation/MMStorage05/Server05/CXC1740305_R9M/storage/corrupt"
	destinoDir := "/mediacao/oi/coll/corrupt/fixa/MGCAlcatel"
	mascara := "GNASD*"
	hexOrigem := "0AA00819"
	hexNew := "0AA00802"

	println("Iniciando tratamento de arquivo corrupt GNASD")

	err := processarArquivos(diretorio, mascara, hexOrigem, hexNew, destinoDir)
	if err != nil {
		fmt.Println("Erro:", err)
	}

}

func substituirHexArquivo(arquivo, hexOrigem, hexNew string) error {
	conteudo, err := os.ReadFile(arquivo)
	if err != nil {
		return fmt.Errorf("erro ao ler o arquivo: %v", err)
	}

	hexOrigemBytes, err := hex.DecodeString(hexOrigem)
	if err != nil {
		return fmt.Errorf("erro ao decodificar string Hex: %v", err)
	}

	hexNewBytes, err := hex.DecodeString(hexNew)
	if err != nil {
		return fmt.Errorf("erro ao decodificar string Hex: %v", err)
	}

	conteudoNovo := strings.ReplaceAll(string(conteudo), string(hexOrigemBytes), string(hexNewBytes))

	err = os.WriteFile(arquivo, []byte(conteudoNovo), 0644)
	if err != nil {
		return fmt.Errorf("erro ao escrever arquivo: %v", err)
	}

	return nil
}

func copiarArquivo(origem, destino string) error {
	origemArquivo, err := os.Open(origem)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo de origem: %v", err)
	}
	defer origemArquivo.Close()

	destinoArquivo, err := os.Create(destino)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo de destino: %v", err)
	}
	defer destinoArquivo.Close()

	_, err = io.Copy(destinoArquivo, origemArquivo)
	if err != nil {
		return fmt.Errorf("erro ao copiar arquivo: %v", err)
	}

	return nil
}

func moverArquivo(origem, destinoDir string) error {
	if _, err := os.Stat(destinoDir); os.IsNotExist(err) {
		return fmt.Errorf("erro: diretório de destino %s não existe", destinoDir)
	}

	nomeArquivo := filepath.Base(origem)
	nomeArquivo = strings.TrimSuffix(nomeArquivo, ":00")
	novoNome := "C-GNASD@" + nomeArquivo
	destino := filepath.Join(destinoDir, novoNome)

	err := copiarArquivo(origem, destino)
	if err != nil {
		return fmt.Errorf("erro ao copiar o arquivo %s para %s: %v", origem, destino, err)
	}

	err = os.Remove(origem)
	if err != nil {
		return fmt.Errorf("erro ao remover arquivo original %s: %v", origem, err)
	}

	fmt.Printf("Movido para: %s\n", destino)
	return nil
}

func processarArquivos(diretorio, mascara, hexOrigem, hexNew, destinoDir string) error {
	padrao := filepath.Join(diretorio, mascara)

	arquivos, err := filepath.Glob(padrao)
	if err != nil {
		return fmt.Errorf("erro ao listar arquivos no diretorio: %v", err)
	}

	if len(arquivos) == 0 {
		fmt.Println("Nenhum arquivo encontrado para tratamento")
		return nil
	}

	for _, arquivo := range arquivos {
		nomeArquivo := filepath.Base(arquivo)
		fmt.Printf("Processando: %s\n", nomeArquivo)

		err := substituirHexArquivo(arquivo, hexOrigem, hexNew)
		if err != nil {
			fmt.Printf("Erro ao processar %s: %v\n", nomeArquivo, err)
		} else {
			fmt.Printf("Tratado: %s\n", nomeArquivo)
		}

		err = moverArquivo(arquivo, destinoDir)
		if err != nil {
			fmt.Printf("Erro ao mover %s: %v\n", arquivo, err)
		}
	}

	return nil
}

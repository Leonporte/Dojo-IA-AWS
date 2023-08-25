package main

import (
	// "context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	// "os"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/textract"
	//"image"
	"log"

	//"github.com/disintegration/imaging"
	//"github.com/nfnt/resize"
	//"golang.org/x/text/width"

	"bytes"

	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
)

func main() {
	http.HandleFunc("/image", imageHandler)
	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Teste imagem!")
	//removeBackgroundHandler(w, r)
	imageTroca("")
}

func removeBackgroundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Teste imagem!")
}

func imageTroca(imageBase64 string) {
	// Decodifica a imagem Base64 para bytes
	imageBytes, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		log.Fatal("Erro ao decodificar imagem Base64:", err)
		return
	}

	config := aws.Config{
		Region:      aws.String("us-east-1"), // Substitua pela região correta
		Credentials: credentials.NewStaticCredentials("AKIAVMATCO2R2W2EQGPI", "gGWo3cFl+LCYzCmTj8i7AkBUzkjpYlyBKPMNWoFk", ""),
	}

	sess := session.Must(session.NewSession(&config))

	// Cria um cliente Textract
	textractClient := textract.New(sess)

	// Prepara a solicitação para o serviço Textract
	input := &textract.DetectDocumentTextInput{
		Document: &textract.Document{
			Bytes: imageBytes,
		},
		// FeatureTypes: []*string{
		// 	aws.String("FORMS"), // Isso instrui o Textract a analisar os campos de formulário e textos
		// },
	}

	// Chama o serviço Textract para analisar o documento
	result, err := textractClient.DetectDocumentText(input)
	if err != nil {
		log.Fatal("Erro ao analisar o documento:", err)
		return
	}

	fmt.Println("Analisa os blocos de texto para encontrar a assinatura do documento")
	// Analisa os blocos de texto para encontrar a assinatura do documento
	for _, block := range result.Blocks {
		// text := block.Text
		// intValue, err := strconv.ParseInt(text, 0, 64)
		// if err != nil {
		// 	fmt.Println("Erro ao converter:", err)
		// 	return
		// }

		// fmt.Println(strconv.FormatInt(intValue), 16)
		if block.Text != nil && strings.Contains(*block.Text, "ASSINATURA") {
			fmt.Println(block)
			boundingBox := block.Geometry.BoundingBox
			fmt.Println(boundingBox)

			// fmt.Println("Assinatura do Documento encontrada em:", *boundingBox, x, y, x+width, y+height)

			// Redimensiona a imagem para o retângulo da assinatura
			//croppedImage := imaging.Crop(imageBytes, image.Rect(x, y, x+width, y+height))

			// Decodifica a imagem base64

			// Converte os bytes decodificados em uma imagem
			srcImage, _, err := image.Decode(bytes.NewReader(imageBytes))
			srcImageConfig, _, err := image.DecodeConfig(bytes.NewReader(imageBytes))
			if err != nil {
				fmt.Println("Erro ao decodificar a imagem:", err)
				return
			}

			x := *boundingBox.Left * float64(srcImageConfig.Width)
			y := *boundingBox.Top * float64(srcImageConfig.Height)
			width := int(x) + int(*boundingBox.Width*float64(srcImageConfig.Width))
			height := int(y) + int(*boundingBox.Height*float64(srcImageConfig.Height))

			//width = width + (width * 0.1)
			x = x - (x * 0.1)
			y = y - (y * 0.1)
			// Define as coordenadas do retângulo de corte (x, y, largura, altura)
			cropRect := image.Rect(int(x), int(y), width, height) //teset

			// Cria uma nova imagem para o corte
			croppedImage := image.NewRGBA(cropRect)

			// Preenche a imagem cortada com uma cor de exemplo (vermelho)

			draw.Draw(croppedImage, croppedImage.Bounds(), srcImage, cropRect.Min, draw.Src)

			// Cria um buffer para armazenar os bytes da imagem cortada em PNG
			var buf bytes.Buffer
			err = png.Encode(&buf, croppedImage)
			if err != nil {
				fmt.Println("Erro ao codificar a imagem cortada:", err)
				return
			}

			// Codifica os bytes da imagem em base64
			base64CroppedImage := base64.StdEncoding.EncodeToString(buf.Bytes())

			file, err := os.Create("output.txt")
			if err != nil {
				fmt.Println("Erro ao criar o arquivo:", err)
				return
			}
			defer file.Close()

			// Escreve a string no arquivo
			_, err = file.WriteString(base64CroppedImage)
			if err != nil {
				fmt.Println("Erro ao escrever no arquivo:", err)
				return
			}

			fmt.Println(srcImageConfig.Width, x, y, width, height)
			fmt.Println(int(x), int(y), int(width), int(height))
			fmt.Println("Imagem cortada salva com sucesso.")
			break
		}
	}

	// // Salva a imagem em uma pasta local
	// err = imaging.Save(croppedImage, "assinatura_recortada.jpg")
	// if err != nil {
	// 	log.Fatal("Erro ao salvar a imagem recortada:", err)
	// }
}

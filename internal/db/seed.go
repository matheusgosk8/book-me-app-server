package db	

import(
	"context"
	"fmt"
	"github.com/matheusgosk8/book-me-server/ent"
)

func SeedDatabase(client *ent.Client) error {
	 ctx := context.Background()
	exists, err := client.Category.Query().Exist(ctx) //Verifica se categoria já existe
	if err != nil {
		return err
	}
	if exists {
		fmt.Println("Banco já semeado. Pulando")
		return nil
	}

	fmt.Println("Semeando categorias...")

	//Criando categoria pai
	limpeza, err := client.Category.Create().SetName("Limpeza").Save(ctx)
	if err != nil {return err}

	reformas, err := client.Category.Create().SetName("Reformas").Save(ctx)
	if err != nil {return err}

	//Subcategorias
	_, err = client.Category.Create().
			SetName("Piscina").	
			SetParent(limpeza).
			Save(ctx)
	if err != nil {return err}

	_, err = client.Category.Create().
			SetName("Pintura").	
			SetParent(reformas).
			Save(ctx)
	if err != nil {return err}

	fmt.Println("Categorias criadas com sucesso!")
	return nil
}
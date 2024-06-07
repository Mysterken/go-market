package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func addProduct(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)

	// 1. Demander à l'utilisateur de saisir les informations du produit
	fmt.Println("Saisir l'identifiant du produit:")
	id, _ := reader.ReadString('\n')

	fmt.Println("Saisir le titre du produit:")
	title, _ := reader.ReadString('\n')

	fmt.Println("Saisir la description du produit:")
	description, _ := reader.ReadString('\n')

	fmt.Println("Saisir le prix du produit:")
	var price float64
	fmt.Scan(&price)

	fmt.Println("Saisir la quantité du produit:")
	var quantity int
	fmt.Scan(&quantity)

	// 2. Insérer les informations du produit dans la table products
	_, err := db.Exec("INSERT INTO products (id, title, description, price, quantity) VALUES (?, ?, ?, ?, ?)", id, title, description, price, quantity)

	if err != nil {
		fmt.Println("Erreur lors de l'insertion du produit")
		return
	}

	// 3. Afficher un message de confirmation
	fmt.Println("Produit ajouté avec succès!")
}

func displayProduct(db *sql.DB, id int) bool {
	// 1. Récupérer les informations du produit
	var title string
	var description string
	var quantity int
	var price float64

	err := db.QueryRow("SELECT title, description, price, quantity FROM products WHERE id = ?", id).Scan(&title, &description, &price, &quantity)

	if err != nil {
		fmt.Println("Produit non trouvé")
		return false
	}

	// 2. Afficher les informations du produit
	fmt.Printf("Titre: %s\n", title)
	fmt.Printf("Description: %s\n", description)
	fmt.Printf("Prix: %.2f\n", price)
	fmt.Printf("Quantité: %d\n", quantity)

	return true
}

func displayProducts(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM products")

	if err != nil {
		fmt.Println("Erreur lors de la récupération des produits")
		return
	}

	defer rows.Close()

	// 2. Afficher les produits
	fmt.Println("Liste des produits:")
	w := tabwriter.NewWriter(os.Stdout, 1, 2, 1, ' ', 0)
	fmt.Fprintln(w, "ID\tTitre\tDescription\tQuantité\tPrix")
	for rows.Next() {
		var id int
		var title string
		var description string
		var price float64
		var quantity int
		rows.Scan(&id, &title, &description, &price, &quantity)

		// trim spaces from title and description
		title = strings.TrimSpace(title)
		description = strings.TrimSpace(description)

		fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%.2f\n", id, title, description, quantity, price)
	}

	w.Flush()
}

func modifyProduct(db *sql.DB) {
	var productChosen bool = false

	for !productChosen {
		fmt.Println("Saisir l'identifiant du produit à modifier:")
		var id int
		fmt.Scan(&id)

		if displayProduct(db, id) {

			// ask if the user wants to modify the product
			fmt.Println("Voulez-vous modifier ce produit? (O/N)")
			var choice string
			fmt.Scan(&choice)

			if choice == "O" || choice == "o" {
				reader := bufio.NewReader(os.Stdin)

				fmt.Println("Saisir le nouveau titre du produit:")
				title, _ := reader.ReadString('\n')

				fmt.Println("Saisir la nouvelle description du produit:")
				description, _ := reader.ReadString('\n')

				fmt.Println("Saisir le nouveau prix du produit:")
				var price float64
				fmt.Scan(&price)

				fmt.Println("Saisir la nouvelle quantité du produit:")
				var quantity int
				fmt.Scan(&quantity)

				// update the product
				_, err := db.Exec("UPDATE products SET title = ?, description = ?, price = ?, quantity = ? WHERE id = ?", title, description, price, quantity, id)

				if err != nil {
					fmt.Println("Erreur lors de la mise à jour du produit")
					return
				}

				fmt.Println("Produit modifié avec succès!")
				productChosen = true
			} else {
				fmt.Println("Opération annulée")
				productChosen = true
			}

		}
	}
}

func deleteProduct(db *sql.DB) {
	var productChosen bool = false

	for !productChosen {
		fmt.Println("Saisir l'identifiant du produit à supprimer:")
		var id int
		fmt.Scan(&id)

		if displayProduct(db, id) {

			// ask if the user wants to delete the product
			fmt.Println("Voulez-vous supprimer ce produit? (O/N)")
			var choice string
			fmt.Scan(&choice)

			if choice == "O" || choice == "o" {
				// delete the product
				_, err := db.Exec("DELETE FROM products WHERE id = ?", id)

				if err != nil {
					fmt.Println("Erreur lors de la suppression du produit")
					return
				}

				fmt.Println("Produit supprimé avec succès!")
				productChosen = true
			} else {
				fmt.Println("Opération annulée")
				productChosen = true
			}

		}

	}
}

func exportProducts(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM products")

	if err != nil {
		fmt.Println("Erreur lors de la récupération des produits")
		return
	}

	defer rows.Close()

	file, err := os.Create("products.csv")
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier")
		return
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	writer.Write([]string{"ID", "Titre", "Description", "Prix", "Quantité"})

	for rows.Next() {
		var id int
		var title string
		var description string
		var price float64
		var quantity int
		rows.Scan(&id, &title, &description, &price, &quantity)

		writer.Write([]string{fmt.Sprintf("%d", id), title, description, fmt.Sprintf("%.2f", price), fmt.Sprintf("%d", quantity)})
	}

	writer.Flush()

	fmt.Println("Exportation des produits terminée")
}

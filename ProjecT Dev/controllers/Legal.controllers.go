package controllers

import (
	models "RedProject/models"
	"fmt"
	"html/template"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

func sendEmail(to, subject, body string) error {
	from := "alexandre.petitfrere@ynov.com"
	addr := os.Getenv("SMTP_ADDR")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	host := os.Getenv("SMTP_HOST")

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s", from, to, subject, body)

	if addr != "" && host != "" {
		auth := smtp.PlainAuth("", smtpUser, smtpPass, host)
		err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
		if err != nil {
			fmt.Printf("Erreur envoi email à %s: %v\n", to, err)
			fmt.Printf("[EMAIL FALLBACK] De: %s\nPour: %s\nSujet: %s\n%s\n", from, to, subject, body)
			return err
		}
		fmt.Printf("Email envoyé à %s depuis %s\n", to, from)
		return nil
	}

	fmt.Printf("[EMAIL FALLBACK] De: %s\nPour: %s\nSujet: %s\n%s\n", from, to, subject, body)
	return nil
}

func LegalHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/red_project/legal/")

	var title string
	var content template.HTML

	switch path {
	case "mentions":
		title = "Mentions légales"
		content = template.HTML(`
			<h2>Éditeur du site</h2>
			<p>YPlaza SAS<br>
			24 rue de la Paix, 75002 Paris<br>
			RCS Paris 850 123 456<br>
			Capital social : 250 000 €<br>
			N° TVA intracommunautaire : FR 12 850123456</p>

			<h2>Directeur de la publication</h2>
			<p>Alexandre Petitfrère, Président-directeur général</p>

			<h2>Hébergement</h2>
			<p>Le site est hébergé par OVH SAS<br>
			2 rue Kellermann, 59100 Roubaix, France</p>

			<h2>Propriété intellectuelle</h2>
			<p>L'ensemble des contenus présents sur le site YPlaza (textes, images, vidéos, logos, marques) est protégé par le droit d'auteur et le droit des marques. Toute reproduction, représentation, modification ou exploitation sans autorisation écrite préalable de YPlaza SAS est interdite.</p>

			<h2>Responsabilité</h2>
			<p>YPlaza s'efforce d'assurer l'exactitude et la mise à jour des informations publiées sur le site. Toutefois, YPlaza ne saurait garantir l'exhaustivité ou l'absence d'erreurs. L'utilisateur reconnaît utiliser ces informations sous sa responsabilité exclusive.</p>
		`)
	case "privacy":
		title = "Politique de confidentialité"
		content = template.HTML(`
			<h2>Collecte des données</h2>
			<p>YPlaza collecte les données personnelles suivantes dans le cadre de ses services : nom, prénom, adresse email, numéro de téléphone, adresse postale. Ces données sont collectées lors de la création de compte, de l'utilisation des services de mise en relation, des enchères et des formulaires de contact.</p>

			<h2>Finalités du traitement</h2>
			<p>Vos données sont utilisées pour :</p>
			<ul>
				<li>La gestion de votre compte et l'accès aux services</li>
				<li>La mise en relation avec les vendeurs et acheteurs</li>
				<li>La participation aux enchères</li>
				<li>L'envoi d'analyses et de prévisions personnalisées</li>
				<li>La réponse à vos demandes via le formulaire de contact</li>
			</ul>

			<h2>Durée de conservation</h2>
			<p>Vos données sont conservées pendant toute la durée de votre compte actif, et jusqu'à 3 ans après votre dernière interaction avec YPlaza.</p>

			<h2>Vos droits</h2>
			<p>Conformément au RGPD, vous disposez d'un droit d'accès, de rectification, de suppression et de portabilité de vos données. Vous pouvez exercer ces droits en nous contactant à contact@yplaza.fr.</p>
		`)
	case "cgu":
		title = "Conditions générales d'utilisation"
		content = template.HTML(`
			<h2>Acceptation des conditions</h2>
			<p>L'utilisation du site YPlaza implique l'acceptation pleine et entière des présentes conditions générales d'utilisation.</p>

			<h2>Services proposés</h2>
			<p>YPlaza propose un service de mise en relation entre vendeurs et acheteurs de biens immobiliers, un système d'enchères en ligne, ainsi qu'un outil d'analyse et de prédiction des prix immobiliers.</p>

			<h2>Obligations de l'utilisateur</h2>
			<p>L'utilisateur s'engage à :</p>
			<ul>
				<li>Fournir des informations exactes lors de son inscription</li>
				<li>Ne pas utiliser le site à des fins frauduleuses</li>
				<li>Respecter les lois et règlements en vigueur</li>
				<li>Ne pas perturber le fonctionnement du site</li>
			</ul>

			<h2>Limitation de responsabilité</h2>
			<p>YPlaza agit comme intermédiaire et ne saurait être tenu responsable des transactions entre utilisateurs, ni des éventuels litiges entre vendeurs et acheteurs.</p>

			<h2>Modification des CGU</h2>
			<p>YPlaza se réserve le droit de modifier les présentes CGU à tout moment. Les utilisateurs seront informés de toute modification substantielle.</p>
		`)
	case "cgv":
		title = "Conditions générales de vente"
		content = template.HTML(`
			<h2>Objet</h2>
			<p>Les présentes conditions générales de vente régissent les transactions réalisées via la plateforme YPlaza entre vendeurs et acheteurs de biens immobiliers.</p>

			<h2>Mise en relation</h2>
			<p>YPlaza met en relation vendeurs et acheteurs. La conclusion de la vente est directement négociée entre les parties. YPlaza n'est pas partie au contrat de vente.</p>

			<h2>Enchères</h2>
			<p>Les enchères sont organisées conformément au règlement spécifique de chaque vente. L'offre la plus élevée au terme de l'enchère emporte le bien, sous réserve de l'accord du vendeur. Les enchères sont fermes et définitives.</p>

			<h2>Frais de service</h2>
			<p>YPlaza perçoit une commission de 3% du prix de vente, à la charge du vendeur, payable au moment de la signature de l'acte authentique.</p>

			<h2>Rétractation</h2>
			<p>Conformément à la loi, l'acheteur non professionnel dispose d'un délai de rétractation de 14 jours à compter de la réservation, sauf pour les biens acquis aux enchères où ce droit ne s'applique pas.</p>
		`)
	default:
		http.NotFound(w, r)
		return
	}

	data := models.LegalPageData{
		Profil:  StructHome.Profil,
		Title:   title,
		Content: content,
	}

	err := temp.ExecuteTemplate(w, "legal", data)
	if err != nil {
		http.Error(w, "Erreur de rendu", http.StatusInternalServerError)
	}
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	data := models.ContactPageData{
		Profil: StructHome.Profil,
	}

	if r.Method == http.MethodPost {
		name := strings.TrimSpace(r.FormValue("name"))
		email := strings.TrimSpace(r.FormValue("email"))
		subject := r.FormValue("subject")
		message := strings.TrimSpace(r.FormValue("message"))

		if name == "" || email == "" || message == "" {
			data.Error = "Tous les champs obligatoires doivent être remplis."
			temp.ExecuteTemplate(w, "contact", data)
			return
		}

		logMsg := "[Contact] De: " + name + " (" + email + ") | Sujet: " + subject + "\n" + message
		if strings.Contains(email, "@") {
			if err := sendEmail("alexandre.petitfrere@ynov.com", "Nouveau message de contact - YPlaza", logMsg); err != nil {
				data.Error = "Erreur lors de l'envoi de votre message. Veuillez réessayer."
				temp.ExecuteTemplate(w, "contact", data)
				return
			}
		} else {
			data.Error = "Adresse email invalide."
			temp.ExecuteTemplate(w, "contact", data)
			return
		}

		data.Success = true
	}

	err := temp.ExecuteTemplate(w, "contact", data)
	if err != nil {
		http.Error(w, "Erreur de rendu", http.StatusInternalServerError)
	}
}

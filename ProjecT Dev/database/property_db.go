package database

import (
	"RedProject/models"
	"database/sql"
	"encoding/json"
)

func GetAllProperties() ([]models.Property, error) {
	rows, err := DB.Query("SELECT id, name, description, price, is_sell, images, property_type, rooms, location, surface FROM properties ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProperties(rows)
}

func GetPropertyByID(id int) (models.Property, error) {
	var p models.Property
	var imagesJSON sql.NullString
	var propType, location sql.NullString
	var rooms sql.NullInt64
	var surface sql.NullFloat64

	err := DB.QueryRow(
		"SELECT id, name, description, price, is_sell, images, property_type, rooms, location, surface FROM properties WHERE id = ?", id,
	).Scan(&p.IdProperty, &p.NameProperty, &p.DescProprety,
		&p.PriceProperty, &p.IsSellProperty, &imagesJSON,
		&propType, &rooms, &location, &surface)
	if err != nil {
		return p, err
	}
	if imagesJSON.Valid {
		json.Unmarshal([]byte(imagesJSON.String), &p.ImgProperty)
	}
	p.Type = propType.String
	p.Rooms = int(rooms.Int64)
	p.Location = location.String
	p.Surface = surface.Float64
	return p, nil
}

func CreateProperty(p models.Property) (int64, error) {
	imagesJSON, _ := json.Marshal(p.ImgProperty)
	res, err := DB.Exec(
		"INSERT INTO properties (name, description, price, is_sell, images, owner_id, property_type, rooms, location, surface) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		p.NameProperty, p.DescProprety, p.PriceProperty, p.IsSellProperty, string(imagesJSON), p.OwnerID,
		p.Type, p.Rooms, p.Location, p.Surface,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func LikeProperty(userID, propertyID int) error {
	_, err := DB.Exec("INSERT IGNORE INTO user_liked_properties (user_id, property_id) VALUES (?, ?)", userID, propertyID)
	return err
}

func UnlikeProperty(userID, propertyID int) error {
	_, err := DB.Exec("DELETE FROM user_liked_properties WHERE user_id = ? AND property_id = ?", userID, propertyID)
	return err
}

func BuyProperty(userID, propertyID int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT IGNORE INTO user_bought_properties (user_id, property_id) VALUES (?, ?)", userID, propertyID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("UPDATE properties SET is_sell = FALSE WHERE id = ?", propertyID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func GetUserUnsoldProperties(userID int) ([]models.Property, error) {
	rows, err := DB.Query(`
		SELECT id, name, description, price, is_sell, images, property_type, rooms, location, surface
		FROM properties WHERE owner_id = ? AND is_sell = TRUE
		ORDER BY id DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProperties(rows)
}

func GetRecommendedProperties(userID int, limit int) ([]models.Property, error) {
	rows, err := DB.Query(`
		SELECT DISTINCT p.id, p.name, p.description, p.price, p.is_sell, p.images,
		       p.property_type, p.rooms, p.location, p.surface
		FROM properties p
		LEFT JOIN user_liked_properties ul ON ul.property_id = p.id AND ul.user_id = ?
		LEFT JOIN user_bought_properties ub ON ub.property_id = p.id AND ub.user_id = ?
		WHERE p.is_sell = TRUE
		AND p.id NOT IN (
			SELECT property_id FROM user_bought_properties WHERE user_id = ?
		)
		AND (
			p.property_type IN (
				SELECT property_type FROM user_liked_properties lp
				JOIN properties pp ON pp.id = lp.property_id
				WHERE lp.user_id = ?
			)
			OR p.location IN (
				SELECT location FROM user_liked_properties lp
				JOIN properties pp ON pp.id = lp.property_id
				WHERE lp.user_id = ?
			)
			OR p.property_type IN (
				SELECT property_type FROM user_bought_properties bp
				JOIN properties pp ON pp.id = bp.property_id
				WHERE bp.user_id = ?
			)
		)
		ORDER BY p.id DESC
		LIMIT ?`, userID, userID, userID, userID, userID, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProperties(rows)
}

func SearchProperties(query string) ([]models.Property, error) {
	rows, err := DB.Query(`
		SELECT id, name, description, price, is_sell, images, property_type, rooms, location, surface
		FROM properties
		WHERE name LIKE ? OR description LIKE ? OR location LIKE ?
		ORDER BY id DESC`, "%"+query+"%", "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProperties(rows)
}

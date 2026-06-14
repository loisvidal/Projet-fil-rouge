package database

import (
	"RedProject/models"
	"database/sql"
	"time"
)

func CreateAuction(propertyID, sellerID int, startPrice, minBidStep float64, endTime time.Time) (int64, error) {
	res, err := DB.Exec(
		`INSERT INTO auctions (property_id, seller_id, start_price, current_price, min_bid_step, start_time, end_time)
		 VALUES (?, ?, ?, ?, ?, NOW(), ?)`,
		propertyID, sellerID, startPrice, startPrice, minBidStep, endTime,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetActiveAuctions() ([]models.Auction, error) {
	rows, err := DB.Query(`
		SELECT a.id, a.property_id, p.name, a.seller_id, a.start_price, a.current_price,
		       a.min_bid_step, a.winner_id, COALESCE(u.name,''), a.start_time, a.end_time, a.is_active,
		       (SELECT COUNT(*) FROM bids WHERE auction_id = a.id)
		FROM auctions a
		JOIN properties p ON p.id = a.property_id
		LEFT JOIN users u ON u.id = a.winner_id
		WHERE a.is_active = TRUE AND a.end_time > NOW()
		ORDER BY a.end_time ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAuctions(rows)
}

func GetAuctionByID(id int) (models.Auction, error) {
	var a models.Auction
	err := DB.QueryRow(`
		SELECT a.id, a.property_id, p.name, a.seller_id, a.start_price, a.current_price,
		       a.min_bid_step, a.winner_id, COALESCE(u.name,''), a.start_time, a.end_time, a.is_active,
		       (SELECT COUNT(*) FROM bids WHERE auction_id = a.id)
		FROM auctions a
		JOIN properties p ON p.id = a.property_id
		LEFT JOIN users u ON u.id = a.winner_id
		WHERE a.id = ?`, id).Scan(
		&a.ID, &a.PropertyID, &a.PropertyName, &a.SellerID, &a.StartPrice, &a.CurrentPrice,
		&a.MinBidStep, &a.WinnerID, &a.WinnerName, &a.StartTime, &a.EndTime, &a.IsActive, &a.BidCount)
	if err != nil {
		return a, err
	}
	return a, nil
}

func PlaceBid(auctionID, userID int, amount float64) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO bids (auction_id, user_id, amount) VALUES (?, ?, ?)", auctionID, userID, amount)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("UPDATE auctions SET current_price = ?, winner_id = ? WHERE id = ?", amount, userID, auctionID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func GetAuctionBids(auctionID int) ([]models.Bid, error) {
	rows, err := DB.Query(`
		SELECT b.id, b.auction_id, b.user_id, COALESCE(u.name,''), b.amount, b.created_at
		FROM bids b
		LEFT JOIN users u ON u.id = b.user_id
		WHERE b.auction_id = ?
		ORDER BY b.created_at DESC`, auctionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []models.Bid
	for rows.Next() {
		var b models.Bid
		if err := rows.Scan(&b.ID, &b.AuctionID, &b.UserID, &b.UserName, &b.Amount, &b.CreatedAt); err != nil {
			return nil, err
		}
		bids = append(bids, b)
	}
	return bids, nil
}

func CloseExpiredAuctions() {
	DB.Exec(`UPDATE auctions SET is_active = FALSE WHERE is_active = TRUE AND end_time <= NOW()`)
}

func scanAuctions(rows *sql.Rows) ([]models.Auction, error) {
	var auctions []models.Auction
	for rows.Next() {
		var a models.Auction
		if err := rows.Scan(&a.ID, &a.PropertyID, &a.PropertyName, &a.SellerID,
			&a.StartPrice, &a.CurrentPrice, &a.MinBidStep,
			&a.WinnerID, &a.WinnerName, &a.StartTime, &a.EndTime, &a.IsActive, &a.BidCount); err != nil {
			return nil, err
		}
		auctions = append(auctions, a)
	}
	return auctions, nil
}

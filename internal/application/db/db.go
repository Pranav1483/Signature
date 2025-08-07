package db

import (
	"fmt"
	"log"
	"signature/internal/entity/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func NewDB(host, port, dbname, user, password string) *DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&tls=false&multiStatements=true&interpolateParams=true", user, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("err connecting to DB: ", err)

	}
	db.AutoMigrate(
		models.Organization{},
	)
	return &DB{
		db: db,
	}
}

func (d *DB) SaveOrganization(orgId, orgName, orgType, url, publicKey, signatureMethod string) error {
	org := models.Organization{
		OrgID:           orgId,
		OrgName:         orgName,
		OrgType:         orgType,
		URL:             url,
		PublicKey:       publicKey,
		SignatureMethod: signatureMethod,
	}
	return d.db.Model(&models.Organization{}).Create(org).Error
}

func (d *DB) GetOrganization(orgId string) (models.Organization, error) {
	var org models.Organization
	err := d.db.Model(&models.Organization{}).Where("org_id = ?", orgId).First(&org).Error
	return org, err
}

func (d *DB) GetOrganizationByType(orgType string) (models.Organization, error) {
	var org models.Organization
	err := d.db.Model(&models.Organization{}).Where("org_type = ?", orgType).First(&org).Error
	return org, err
}

func (d *DB) GetOrganizationByURL(url string) (models.Organization, error) {
	var org models.Organization
	err := d.db.Model(&models.Organization{}).Where("hostname LIKE ?", url).First(&org).Error
	return org, err
}

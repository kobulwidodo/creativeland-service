package umkm

import (
	"database/sql"
	"database/sql/driver"
	"go-clean/src/business/entity"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_umkm_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "INSERT INTO"
	query := regexp.QuoteMeta(querySql)

	mockUmkm := entity.Umkm{
		Name: "umkm 1",
	}

	type args struct {
		umkm entity.Umkm
	}
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        entity.Umkm
		wantErr     bool
	}{
		{
			name: "failed to create umkm",
			args: args{
				umkm: mockUmkm,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			want:    mockUmkm,
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				umkm: mockUmkm,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit()
				sqlMock.ExpectationsWereMet()
				return sqlServer, err
			},
			want:    mockUmkm,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlServer, err := tt.prepSqlMock()
			if err != nil {
				t.Error(err)
			}
			defer sqlServer.Close()

			sqlClient, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlServer,
				SkipInitializeWithVersion: true,
			}))
			if err != nil {
				t.Error(err)
			}

			u := Init(sqlClient)
			_, err = u.Create(tt.args.umkm)
			if (err != nil) != tt.wantErr {
				t.Errorf("umkm.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_umkm_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "SELECT * FROM `umkms` WHERE `umkms`.`deleted_at` IS NULL"
	query := regexp.QuoteMeta(querySql)

	mockParam := entity.UmkmParam{}

	type args struct {
		param entity.UmkmParam
	}
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        []entity.Umkm
		wantErr     bool
	}{
		{
			name: "failed to exec query",
			args: args{
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			want:    []entity.Umkm{},
			wantErr: true,
		},
		{
			name: "all ok",
			args: args{
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"name"})
				row.AddRow("umkm 1")
				sqlMock.ExpectQuery(query).WillReturnRows(row)
				return sqlServer, err
			},
			want: []entity.Umkm{
				{
					Name: "umkm 1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlServer, err := tt.prepSqlMock()
			if err != nil {
				t.Error(err)
			}
			defer sqlServer.Close()

			sqlClient, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlServer,
				SkipInitializeWithVersion: true,
			}))
			if err != nil {
				t.Error(err)
			}

			u := Init(sqlClient)
			got, err := u.GetList(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("umkm.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_umkm_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "SELECT * FROM `umkms` WHERE `umkms`.`deleted_at` IS NULL"
	query := regexp.QuoteMeta(querySql)

	mockParam := entity.UmkmParam{}

	type args struct {
		param entity.UmkmParam
	}
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        entity.Umkm
		wantErr     bool
	}{
		{
			name: "failed to exec query",
			args: args{
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			want:    entity.Umkm{},
			wantErr: true,
		},
		{
			name: "all ok",
			args: args{
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"name"})
				row.AddRow("umkm 1")
				sqlMock.ExpectQuery(query).WillReturnRows(row)
				return sqlServer, err
			},
			want: entity.Umkm{
				Name: "umkm 1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlServer, err := tt.prepSqlMock()
			if err != nil {
				t.Error(err)
			}
			defer sqlServer.Close()

			sqlClient, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlServer,
				SkipInitializeWithVersion: true,
			}))
			if err != nil {
				t.Error(err)
			}

			u := Init(sqlClient)
			got, err := u.Get(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("umkm.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_umkm_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "UPDATE"
	query := regexp.QuoteMeta(querySql)

	selectParam := entity.UmkmParam{
		ID: 1,
	}

	updateParam := entity.UpdateUmkmParam{
		Name: "umkm 1 baru",
	}

	type args struct {
		selectParam entity.UmkmParam
		updateParam entity.UpdateUmkmParam
	}
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		wantErr     bool
	}{
		{
			name: "failed to exec query",
			args: args{
				selectParam: selectParam,
				updateParam: updateParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			wantErr: true,
		},
		{
			name: "all ok",
			args: args{
				selectParam: selectParam,
				updateParam: updateParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(driver.RowsAffected(1))
				sqlMock.ExpectCommit()
				sqlMock.ExpectationsWereMet()
				return sqlServer, err
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlServer, err := tt.prepSqlMock()
			if err != nil {
				t.Error(err)
			}
			defer sqlServer.Close()

			sqlClient, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlServer,
				SkipInitializeWithVersion: true,
			}))
			if err != nil {
				t.Error(err)
			}

			u := Init(sqlClient)
			err = u.Update(tt.args.selectParam, tt.args.updateParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("umkm.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_umkm_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "UPDATE"
	query := regexp.QuoteMeta(querySql)

	mockParam := entity.UmkmParam{
		ID: 1,
	}

	type args struct {
		param entity.UmkmParam
	}
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		wantErr     bool
	}{
		{
			name: "failed to exec query",
			args: args{
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			wantErr: true,
		},
		{
			name: "all ok",
			args: args{
				param: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(driver.RowsAffected(1))
				sqlMock.ExpectCommit()
				sqlMock.ExpectationsWereMet()
				return sqlServer, err
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlServer, err := tt.prepSqlMock()
			if err != nil {
				t.Error(err)
			}
			defer sqlServer.Close()

			sqlClient, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlServer,
				SkipInitializeWithVersion: true,
			}))
			if err != nil {
				t.Error(err)
			}

			u := Init(sqlClient)
			err = u.Delete(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("umkm.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

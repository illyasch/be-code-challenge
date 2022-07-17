package handlers_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ardanlabs/conf/v3"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/illyasch/be-code-challenge/cmd/challenge/handlers"
	"github.com/illyasch/be-code-challenge/pkg/business/calc"
	"github.com/illyasch/be-code-challenge/pkg/data/database"
	"github.com/illyasch/be-code-challenge/pkg/sys/logger"
)

var (
	postgresDB *sqlx.DB
	stdLgr     *zap.SugaredLogger
)

func TestMain(m *testing.M) {
	var err error
	stdLgr, err = logger.New("challenge")
	if err != nil {
		log.Fatal(err)
	}

	cfg := struct {
		conf.Version
		DB struct {
			User         string `conf:"default:test"`
			Password     string `conf:"default:test,mask"`
			Host         string `conf:"default:localhost"`
			Name         string `conf:"default:eth"`
			MaxIdleConns int    `conf:"default:0"`
			MaxOpenConns int    `conf:"default:0"`
			DisableTLS   bool   `conf:"default:true"`
		}
	}{
		Version: conf.Version{
			Build: "test",
			Desc:  "Copyright Ilya Scheblanov",
		},
	}

	const prefix = "CHALLENGE"
	_, err = conf.Parse(prefix, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cfg)

	postgresDB, err = database.Open(database.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	})
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestAPIConfig_handleHourly(t *testing.T) {
	t.Parallel()
	cfg := handlers.APIConfig{
		Log: stdLgr,
		DB:  postgresDB,
	}

	t.Run("successful fees calculation", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest(http.MethodGet, "/hourly", nil)
		w := httptest.NewRecorder()

		cfg.Router().ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var got []calc.Hour
		err := json.NewDecoder(w.Body).Decode(&got)
		require.NoError(t, err)

		exp := []calc.Hour{
			{1599436800, 17781937815.707344},
			{1599440400, 25796173158.88589},
			{1599444000, 34821055861.44104},
			{1599447600, 29814493424.40487},
			{1599451200, 27821774201.37403},
			{1599454800, 25575311595.65763},
			{1599458400, 33138772595.681362},
			{1599462000, 35671504748.405235},
			{1599465600, 29861742077.137997},
			{1599469200, 31870528305.547283},
			{1599472800, 29492777739.82145},
			{1599476400, 28476895183.216103},
			{1599480000, 31458479835.443005},
			{1599483600, 36114881483.45387},
			{1599487200, 39990571952.80128},
			{1599490800, 32351461366.072742},
			{1599494400, 35702769826.03566},
			{1599498000, 28225350833.576153},
			{1599501600, 23619974534.4661},
			{1599505200, 20792555662.410378},
			{1599508800, 20324756156.83894},
			{1599512400, 18641503209.089615},
			{1599516000, 16778240397.619787},
			{1599519600, 17399949324.595974},
		}

		assert.Equal(t, exp, got)
	})
}

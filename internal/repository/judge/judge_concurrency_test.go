package judge_test

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"godance/internal/domain"
	"godance/internal/models"
	judgerepo "godance/internal/repository/judge"
	types "godance/internal/type"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const judgeLimit = 5

func testDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := getenvOrSkip(t, "TEST_DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	require.NoError(t, err)

	require.NoError(t, db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Competition{},
		&models.JudgesCompetition{},
	))

	require.NoError(t, db.Exec("TRUNCATE judges_competitions RESTART IDENTITY CASCADE").Error)

	return db
}

func getenvOrSkip(t *testing.T, key string) string {
	t.Helper()
	v := osGetenv(key)
	if v == "" {
		t.Skipf("%s не задана — тест требует реальной БД PostgreSQL", key)
	}
	return v
}

// seedOrganizer создаёт пользователя-организатора и возвращает его ID.
func seedOrganizer(t *testing.T, db *gorm.DB) uint {
	t.Helper()
	u := models.User{
		Email:        uniqueEmail("organizer"),
		PasswordHash: "x",
		Phone:        uniquePhone(),
		Role:         string(types.UserRoleOrganizer),
	}
	require.NoError(t, db.Create(&u).Error)
	return u.ID
}

func seedScenario(t *testing.T, db *gorm.DB, organizerID uint, n int) (uint, []uint) {
	t.Helper()

	comp := models.Competition{
		OrganizerID: organizerID,
		Title:       "concurrency test",
		Status:      string(types.CompetitionStatusOpen),
	}
	require.NoError(t, db.Create(&comp).Error)

	judgeIDs := make([]uint, n)
	for i := 0; i < n; i++ {
		u := models.User{
			Email:        uniqueEmail(fmt.Sprintf("judge%d", i)),
			PasswordHash: "x",
			Phone:        uniquePhone(),
			Role:         string(types.UserRoleJudge),
		}
		require.NoError(t, db.Create(&u).Error)
		require.NoError(t, db.Create(&models.Profile{UserID: u.ID}).Error)
		judgeIDs[i] = u.ID
	}
	return comp.ID, judgeIDs
}

func TestAssign_ConcurrentDistinctJudges_RespectsLimit(t *testing.T) {
	db := testDB(t)
	repo := judgerepo.NewRepository(db)
	organizerID := seedOrganizer(t, db)

	const attempts = 8
	compID, judgeIDs := seedScenario(t, db, organizerID, attempts)

	var success, limited int64
	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func(jid uint) {
			defer wg.Done()
			<-start
			_, err := repo.Assign(compID, jid, organizerID, judgeLimit)
			switch {
			case err == nil:
				atomic.AddInt64(&success, 1)
			case errors.Is(err, domain.ErrJudgeLimitReached):
				atomic.AddInt64(&limited, 1)
			default:
				t.Errorf("неожиданная ошибка: %v", err)
			}
		}(judgeIDs[i])
	}

	close(start)
	wg.Wait()

	require.EqualValues(t, judgeLimit, success,
		"должно пройти ровно %d назначений", judgeLimit)
	require.EqualValues(t, attempts-judgeLimit, limited,
		"остальные попытки отклонены по лимиту")

	var stored int64
	require.NoError(t, db.Model(&models.JudgesCompetition{}).
		Where("competition_id = ?", compID).Count(&stored).Error)
	require.EqualValues(t, judgeLimit, stored,
		"в БД сохранено не более лимита записей")
}

func TestAssign_ConcurrentSameJudge_NoDuplicate(t *testing.T) {
	db := testDB(t)
	repo := judgerepo.NewRepository(db)
	organizerID := seedOrganizer(t, db)

	const attempts = 6
	compID, judgeIDs := seedScenario(t, db, organizerID, 1)
	jid := judgeIDs[0]

	var success, dup int64
	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, err := repo.Assign(compID, jid, organizerID, judgeLimit)
			switch {
			case err == nil:
				atomic.AddInt64(&success, 1)
			case errors.Is(err, domain.ErrAlreadyAssigned):
				atomic.AddInt64(&dup, 1)
			default:
				t.Errorf("неожиданная ошибка: %v", err)
			}
		}()
	}

	close(start)
	wg.Wait()

	require.EqualValues(t, 1, success, "успешным может быть только одно назначение")
	require.EqualValues(t, attempts-1, dup, "остальные отклонены как дубликаты")

	var stored int64
	require.NoError(t, db.Model(&models.JudgesCompetition{}).
		Where("competition_id = ?", compID).Count(&stored).Error)
	require.EqualValues(t, 1, stored, "в БД ровно одна запись о назначении")
}

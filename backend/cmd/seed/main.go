package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

const (
	seedDomain = "@seed.local"
	seedPass   = "secret123"
	batchSize  = 2000
)

var (
	flagUsers        = flag.Int("users", 100000, "number of users to create (85% seekers, 15% employers)")
	flagCompanies    = flag.Int("companies", 8000, "number of companies (subset of employers)")
	flagJobs         = flag.Int("jobs", 25000, "number of jobs")
	flagApplications = flag.Int("applications", 80000, "number of applications")
	flagSavedJobs    = flag.Int("saved-jobs", 30000, "number of saved jobs")
	flagSavedSearch  = flag.Int("saved-searches", 10000, "number of saved searches")
	flagClean        = flag.Bool("clean", false, "delete previously seeded @seed.local users first")
	flagUploads      = flag.String("uploads", "./uploads", "path to uploads directory")
)

type seeder struct {
	pool       *pgxpool.Pool
	r          *rand.Rand
	passHash   string
	resumePath string
	resumeURL  string
	logoPaths  []string
	logoURLs   []string

	seekerIDs   []uuid.UUID
	employerIDs []uuid.UUID
	companyIDs  []uuid.UUID
	jobIDs      []uuid.UUID // published jobs only
}

func main() {
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Println("no .env found, using environment")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	s := &seeder{
		pool: pool,
		r:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	// Compute password hash once — reused for every seeded user.
	hash, err := bcrypt.GenerateFromPassword([]byte(seedPass), 12)
	if err != nil {
		log.Fatalf("bcrypt: %v", err)
	}
	s.passHash = string(hash)
	fmt.Println("✓ password hash computed (shared across all seeded users)")

	// Ensure dummy files exist.
	resumePath, logoPaths, err := ensureAssets(*flagUploads)
	if err != nil {
		log.Fatalf("assets: %v", err)
	}
	s.resumePath = resumePath
	s.logoPaths = logoPaths
	s.resumeURL = "/uploads/resumes/" + filenameOf(resumePath)
	for _, p := range logoPaths {
		s.logoURLs = append(s.logoURLs, "/uploads/logos/"+filenameOf(p))
	}
	fmt.Printf("✓ placeholder assets written (%d logos, 1 resume)\n", len(logoPaths))

	if *flagClean {
		fmt.Println("⚠  cleaning existing @seed.local data…")
		if err := s.clean(ctx); err != nil {
			log.Fatalf("clean: %v", err)
		}
	}

	total := time.Now()

	fmt.Printf("→ seeding %d users…\n", *flagUsers)
	if err := s.seedUsers(ctx, *flagUsers); err != nil {
		log.Fatalf("users: %v", err)
	}
	fmt.Printf("  ✓ %d seekers, %d employers\n", len(s.seekerIDs), len(s.employerIDs))

	fmt.Printf("→ seeding companies (%d)…\n", *flagCompanies)
	if err := s.seedCompanies(ctx, *flagCompanies); err != nil {
		log.Fatalf("companies: %v", err)
	}
	fmt.Printf("  ✓ %d companies\n", len(s.companyIDs))

	fmt.Printf("→ seeding jobs (%d)…\n", *flagJobs)
	if err := s.seedJobs(ctx, *flagJobs); err != nil {
		log.Fatalf("jobs: %v", err)
	}
	fmt.Printf("  ✓ %d published jobs, %d total\n", len(s.jobIDs), *flagJobs)

	fmt.Printf("→ seeding applications (%d)…\n", *flagApplications)
	if err := s.seedApplications(ctx, *flagApplications); err != nil {
		log.Fatalf("applications: %v", err)
	}

	fmt.Printf("→ seeding saved jobs (%d)…\n", *flagSavedJobs)
	if err := s.seedSavedJobs(ctx, *flagSavedJobs); err != nil {
		log.Fatalf("saved jobs: %v", err)
	}

	fmt.Printf("→ seeding saved searches (%d)…\n", *flagSavedSearch)
	if err := s.seedSavedSearches(ctx, *flagSavedSearch); err != nil {
		log.Fatalf("saved searches: %v", err)
	}

	// Stable test accounts.
	if err := s.createTestAccounts(ctx); err != nil {
		log.Fatalf("test accounts: %v", err)
	}

	fmt.Printf("\n✅ done in %s\n\n", time.Since(total).Round(time.Second))
	fmt.Println("Test accounts (password: " + seedPass + ")")
	fmt.Println("  admin1" + seedDomain)
	fmt.Println("  seeker1" + seedDomain + "   (completed profile, resume, applications, saves)")
	fmt.Println("  employer1" + seedDomain + " (company, published jobs, applications)")
	fmt.Println()
	fmt.Printf("Total seeded: %d users, %d companies, %d jobs\n",
		len(s.seekerIDs)+len(s.employerIDs), len(s.companyIDs), *flagJobs)
}

func filenameOf(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

func (s *seeder) clean(ctx context.Context) error {
	// Deleting users cascades to profiles, companies, jobs, applications, saved_*, notifications.
	// Job rows cascade too. Everything seeded has email LIKE '%@seed.local'.
	_, err := s.pool.Exec(ctx, `DELETE FROM users WHERE email LIKE '%@seed.local'`)
	return err
}

// copyFrom is a thin helper around pgx.CopyFrom.
func (s *seeder) copyFrom(ctx context.Context, table string, cols []string, rows [][]any) (int64, error) {
	return s.pool.CopyFrom(
		ctx,
		pgx.Identifier{table},
		cols,
		pgx.CopyFromRows(rows),
	)
}

package main

import (
  "database/sql"
  "encoding/json"
  "log"
  "net/http"
  "os"
  "path/filepath"

  _ "github.com/lib/pq"
)

type Item struct {
  ID   int    `json:"id"`
  Name string `json:"name"`
}

func env(key, def string) string {
  if v := os.Getenv(key); v != "" {
    return v
  }
  return def
}

func main() {
  dsn := "host=" + env("DB_HOST", "db") +
    " port=" + env("DB_PORT", "5432") +
    " user=" + env("DB_USER", "") +
    " password=" + env("DB_PASSWORD", "") +
    " dbname=" + env("DB_NAME", "") +
    " sslmode=disable"

  db, err := sql.Open("postgres", dsn)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  staticDir := env("STATIC_DIR", "/app/static")

  http.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, name FROM items ORDER BY id")
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
    defer rows.Close()

    items := []Item{}
    for rows.Next() {
      var it Item
      if err := rows.Scan(&it.ID, &it.Name); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
      }
      items = append(items, it)
    }
    if err := rows.Err(); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(items)
  })

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/" {
      http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
      return
    }
    http.FileServer(http.Dir(staticDir)).ServeHTTP(w, r)
  })

  addr := ":" + env("HTTP_PORT", "8080")
  log.Printf("listening on %s", addr)
  if err := http.ListenAndServe(addr, nil); err != nil {
    log.Fatal(err)
  }
}
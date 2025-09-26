package main

import (
    "database/sql"
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/lib/pq"
    _ "github.com/lib/pq"
)

// Структуры запросов
type SensorIDsRequest struct {
    SensorIDs  []uuid.UUID `json:"sensorIds" binding:"required"`
    DateFrom   time.Time   `json:"dateFrom" binding:"required"`
    DateTo     time.Time   `json:"dateTo" binding:"required"`
}

type LocationRequest struct {
    Location uuid.UUID `json:"location" binding:"required"`
    DateFrom time.Time `json:"dateFrom" binding:"required"`
    DateTo   time.Time `json:"dateTo" binding:"required"`
}

// Структуры ответов
type SensorInfo struct {
    SensorID   uuid.UUID `json:"sensorId"`
    SensorName string    `json:"sensorName"`
    Value      float64   `json:"value"`
    Timestamp  time.Time `json:"timestamp"`
}

type SensorsInfo struct {
    Data []SensorInfo `json:"data"`
}

type SensorsLocationInfo struct {
    LocationID uuid.UUID    `json:"locationId"`
    Sensors    []SensorInfo `json:"sensors"`
}

type SensorStatus struct {
    SensorID    uuid.UUID `json:"sensorId"`
    SensorName  string    `json:"sensorName"`
    LastValue   float64   `json:"lastValue"`
    LastReading time.Time `json:"lastReading"`
    IsActive    bool      `json:"isActive"`
}

type SensorsStatus struct {
    Statuses []SensorStatus `json:"statuses"`
}

var db *sql.DB

func main() {
    // Подключение к БД
    connStr := "host=postgres port=5432 user=postgres password=postgres dbname=telemetry sslmode=disable"
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Проверка подключения
    if err = db.Ping(); err != nil {
        log.Fatal("Failed to ping database:", err)
    }

    router := gin.Default()

    // CORS middleware
    router.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    })

    // API endpoints
    router.POST("/api/v1/sensors/info", getSensorsInfo)
    router.POST("/api/v1/sensors/location", getSensorsLocationInfo)
    router.POST("/api/v1/sensors/status", getSensorsStatus)

    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now()})
    })

    // Swagger UI
    router.Static("/swagger", "./swagger")
    router.GET("/", func(c *gin.Context) {
        c.Redirect(302, "/swagger/index.html")
    })

    log.Println("Server starting on :8080")
    log.Fatal(router.Run(":8080"))
}

func getSensorsInfo(c *gin.Context) {
    var req SensorIDsRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
        return
    }

    query := `
        SELECT s.id, s.name, td.value, td.timestamp 
        FROM telemetry_data td
        JOIN sensors s ON td.sensor_id = s.id
        WHERE s.id = ANY($1) AND td.timestamp BETWEEN $2 AND $3
        ORDER BY td.timestamp DESC
    `

    // Преобразуем UUID в строки для pq.Array
    sensorIDs := make([]string, len(req.SensorIDs))
    for i, id := range req.SensorIDs {
        sensorIDs[i] = id.String()
    }

    rows, err := db.Query(query, pq.Array(sensorIDs), req.DateFrom, req.DateTo)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
        return
    }
    defer rows.Close()

    var result SensorsInfo
    for rows.Next() {
        var info SensorInfo
        var sensorID string
        if err := rows.Scan(&sensorID, &info.SensorName, &info.Value, &info.Timestamp); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
            return
        }
        info.SensorID, _ = uuid.Parse(sensorID)
        result.Data = append(result.Data, info)
    }

    c.JSON(http.StatusOK, result)
}

func getSensorsLocationInfo(c *gin.Context) {
    var req LocationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
        return
    }

    query := `
        SELECT s.id, s.name, td.value, td.timestamp 
        FROM telemetry_data td
        JOIN sensors s ON td.sensor_id = s.id
        WHERE s.location_id = $1 AND td.timestamp BETWEEN $2 AND $3
        ORDER BY td.timestamp DESC
    `

    rows, err := db.Query(query, req.Location.String(), req.DateFrom, req.DateTo)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
        return
    }
    defer rows.Close()

    var result SensorsLocationInfo
    result.LocationID = req.Location

    for rows.Next() {
        var info SensorInfo
        var sensorID string
        if err := rows.Scan(&sensorID, &info.SensorName, &info.Value, &info.Timestamp); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
            return
        }
        info.SensorID, _ = uuid.Parse(sensorID)
        result.Sensors = append(result.Sensors, info)
    }

    c.JSON(http.StatusOK, result)
}

func getSensorsStatus(c *gin.Context) {
    var req SensorIDsRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
        return
    }

    query := `
        SELECT s.id, s.name, td.value, td.timestamp,
               (EXTRACT(EPOCH FROM (NOW() - td.timestamp)) < 3600) as is_active
        FROM sensors s
        LEFT JOIN LATERAL (
            SELECT value, timestamp 
            FROM telemetry_data 
            WHERE sensor_id = s.id AND timestamp BETWEEN $2 AND $3
            ORDER BY timestamp DESC 
            LIMIT 1
        ) td ON true
        WHERE s.id = ANY($1)
    `

    sensorIDs := make([]string, len(req.SensorIDs))
    for i, id := range req.SensorIDs {
        sensorIDs[i] = id.String()
    }

    rows, err := db.Query(query, pq.Array(sensorIDs), req.DateFrom, req.DateTo)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
        return
    }
    defer rows.Close()

    var result SensorsStatus
    for rows.Next() {
        var status SensorStatus
        var sensorID string
        var value sql.NullFloat64
        var timestamp sql.NullTime
        
        if err := rows.Scan(&sensorID, &status.SensorName, &value, &timestamp, &status.IsActive); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Scan error: " + err.Error()})
            return
        }
        
        status.SensorID, _ = uuid.Parse(sensorID)
        if value.Valid {
            status.LastValue = value.Float64
        }
        if timestamp.Valid {
            status.LastReading = timestamp.Time
        } else {
            status.IsActive = false
        }

        result.Statuses = append(result.Statuses, status)
    }

    c.JSON(http.StatusOK, result)
}
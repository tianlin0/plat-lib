package task

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/tianlin0/plat-lib/conn"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

var (
	commitTask Task
	jobIDMap   = make(map[int64]*mysqlJob)
	testDB     *sql.DB
	testClient *redis.Client
)

const (
	queryJob = `SELECT id, name, task_id, task_name, status, payload, failed_times, commit_id, created_at, updated_at FROM %s WHERE task_id=?;`
)

var mysqlConn = &conn.Connect{
	Host:     "127.0.0.1",
	Port:     "3306",
	Username: "root",
	Password: "root",
	Database: "db_task_test",
}

// JobFunc example
func MyJob(ctx context.Context, id, fts, cID int64, payload []byte) error {

	fmt.Println("MyJob assigned job: ", id, string(payload))

	// then, do something
	fmt.Println("do something ...")

	// return if error with succ=false
	if id%99 == 0 {
		return fmt.Errorf("mock error")
	}

	return nil

}

func TestCreateDB(t *testing.T) {

	mytask, err := NewTask(context.Background(),
		"my_test_task",
		mysqlConn,
		testClient,
		SetJobFunc(MyJob),
		SetMaxFailedTimes(3),
		SetParallelNum(10),
		SetHandleZombieJobsDstStatus(JobStatusTodo),
		SetHandleZombieJobsPeriod(time.Second*5),
		SetHandleZombieJobsTimeout(time.Hour*24),
		SetHandleZombieJobsLimit(100),
		SetTableName("task_test_table"),
		SetInitDB(true),
	)

	fmt.Println(mytask, err)

	return
}

//func TestCreateDB(t *testing.T) error {
//
//	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", "?parseTime=true&charset=utf8"}, "")
//	db, err := sql.Open("mysql", path)
//	if err != nil {
//		return err
//	}
//	defer db.Close()
//
//	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s;", dbName))
//	if err != nil {
//		return err
//	}
//
//	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", dbName))
//	if err != nil {
//		return err
//	}
//
//	_, err = db.Exec("USE " + dbName)
//	if err != nil {
//		return err
//	}
//
//	sql := `CREATE TABLE %s (
//		id bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键, 无业务意义',
//		name varchar(255) NOT NULL COMMENT 'job name, 便于查询',
//		task_id varchar(255) NOT NULL COMMENT 'task id',
//		task_name varchar(255) NOT NULL COMMENT 'task name',
//		status varchar(255) NOT NULL COMMENT 'status',
//		payload blob NOT NULL COMMENT 'payload',
//		failed_times varchar(64) NOT NULL DEFAULT '0' COMMENT 'failed_times',
//		commit_id varchar(255) NOT NULL DEFAULT '0' COMMENT 'commit id: 乐观锁',
//		created_at varchar(64) NOT NULL DEFAULT '0' COMMENT '创建时间',
//		updated_at varchar(64) NOT NULL DEFAULT '0' COMMENT '更新时间',
//		PRIMARY KEY (id),
//		KEY idx_status (status(32)),
//		KEY idx_task_id (task_id(32)),
//		KEY idx_updated_at (updated_at(64)),
//		KEY idx_name (name(32))
//	) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='common task job';`
//
//	_, err = db.Exec(fmt.Sprintf(sql, "task_test_table"))
//	if err != nil {
//		return err
//	}
//
//	_, err = db.Exec(fmt.Sprintf(sql, "task_test_table_with_commit"))
//	if err != nil {
//		return err
//	}
//
//	insertZombieSql := `INSERT INTO %s (name, task_id, task_name, status, payload, failed_times, commit_id, created_at, updated_at)
//	VALUES (?,?,?,?,?,?,?,?,?);`
//	// mock zombie job
//	for i := 0; i < 400; i++ {
//		_, err = db.Exec(fmt.Sprintf(insertZombieSql, "task_test_table"),
//			"job_rand_name_for_ut",
//			"my_test_task",
//			"my_test_task",
//			"todo",
//			[]byte("hello"),
//			0,
//			0,
//			time.Now().AddDate(0, 0, -3).Unix(),
//			time.Now().AddDate(0, 0, -3).Unix(),
//		)
//		if err != nil {
//			return err
//		}
//	}
//
//	// mock zombie job
//	for i := 0; i < 400; i++ {
//		_, err = db.Exec(fmt.Sprintf(insertZombieSql, "task_test_table_with_commit"),
//			"job_rand_name_for_ut",
//			"my_test_task_with_commit",
//			"my_test_task_with_commit",
//			"todo",
//			[]byte("hello"),
//			0,
//			0,
//			time.Now().AddDate(0, 0, -3).Unix(),
//			time.Now().AddDate(0, 0, -3).Unix(),
//		)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

//func testInit() (*sql.DB, *redis.Client, error) {
//	client := redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		Password: "",
//		DB:       0,
//	})
//
//	pong, err := client.Ping(context.Background()).Result()
//	fmt.Println(pong, err)
//	if err != nil {
//		return nil, nil, fmt.Errorf("ping redis error:%v", err)
//	}
//
//	err = createDB()
//	if err != nil {
//		return nil, nil, fmt.Errorf("createDB error:%v", err)
//	}
//
//	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?parseTime=true&charset=utf8"}, "")
//	db, err := sql.Open("mysql", path)
//	if err != nil {
//		return nil, nil, fmt.Errorf("open path:%s, error:%v", path, err)
//	}
//
//	db.SetConnMaxLifetime(100)
//	db.SetMaxIdleConns(10)
//	if err := db.Ping(); err != nil {
//		return nil, nil, fmt.Errorf("ping db path:%s, error:%v", path, err)
//	}
//
//	testDB, testClient = db, client
//
//	return db, client, nil
//}
//
//func TestMustNewTask(t *testing.T) {
//	testDB, testClient, err := testInit()
//
//	mytask := MustNewTask(context.Background(),
//		"my_test_task",
//		testDB,
//		testClient,
//		SetJobFunc(MyJob),
//		SetMaxFailedTimes(3),
//		SetParallelNum(10),
//		SetHandleZombieJobsDstStatus(JobStatusTodo),
//		SetHandleZombieJobsPeriod(time.Second*5),
//		SetHandleZombieJobsTimeout(time.Hour*24),
//		SetHandleZombieJobsLimit(100),
//		SetTableName("task_test_table"),
//		SetInitDB(false),
//	)
//
//	// run consumer 2 to mock another process
//	MustNewTask(context.Background(),
//		"my_test_task",
//		testDB,
//		testClient,
//		SetJobFunc(MyJob),
//		SetMaxFailedTimes(3),
//		SetParallelNum(10),
//		SetHandleZombieJobsDstStatus(JobStatusTodo),
//		SetHandleZombieJobsPeriod(time.Second*5),
//		SetHandleZombieJobsTimeout(time.Hour*24),
//		SetHandleZombieJobsLimit(100),
//		SetTableName("task_test_table"),
//		SetInitDB(false),
//	)
//
//	// run consumer 3 to mock another process
//	MustNewTask(context.Background(),
//		"my_test_task",
//		testDB,
//		testClient,
//		SetJobFunc(MyJob),
//		SetMaxFailedTimes(3),
//		SetParallelNum(10),
//		SetHandleZombieJobsDstStatus(JobStatusTodo),
//		SetHandleZombieJobsPeriod(time.Second*5),
//		SetHandleZombieJobsTimeout(time.Hour*24),
//		SetHandleZombieJobsLimit(100),
//		SetTableName("task_test_table"),
//		SetInitDB(false),
//	)
//
//	// mock push data
//	totalJobs := 300
//	for i := 0; i < totalJobs; i++ {
//		mytask.Push([]byte("hello"))
//	}
//
//	time.Sleep(time.Second * 30)
//
//	// check result
//	list := make([]*mysqlJob, 0)
//	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
//	defer cancel()
//	rows, err := testDB.QueryContext(ctx, fmt.Sprintf(queryJob, "task_test_table"), "my_test_task")
//
//	fmt.Println(err)
//
//	for rows.Next() {
//		job := new(mysqlJob)
//		// id, name, task_id, task_name, status, payload, failed_times, created_at, updated_at
//		err = rows.Scan(&job.Id, &job.Name, &job.TaskID, &job.TaskName, &job.Status, &job.Payload, &job.FailedTimes, &job.CommitID, &job.CreatedAt, &job.UpdatedAt)
//
//		list = append(list, job)
//	}
//	rows.Close()
//
//	failed := 0
//	todo := 0
//	done := 0
//	for _, job := range list {
//		switch JobStatus(job.Status) {
//		case JobStatusDone:
//			done++
//		case JobStatusFailed:
//			failed++
//		case JobStatusTodo:
//			todo++
//		}
//	}
//
//}
//
//func TestMustNewTaskWithCommit(t *testing.T) {
//	commitTask = MustNewTask(context.Background(),
//		"my_test_task_with_commit",
//		testDB,
//		testClient,
//		SetJobFunc(MyJobWithCommit),
//		SetParallelNum(10),
//		// SetHandleZombieJobsDstStatus(JobStatusTodo),// defualt: zombie
//		SetHandleZombieJobsPeriod(time.Second*5),
//		SetHandleZombieJobsTimeout(time.Hour*24),
//		SetHandleZombieJobsLimit(100),
//		SetTableName("task_test_table_with_commit"),
//		SetInitDB(false),
//		SetAutoCommit(false),
//	)
//
//	// run consumer 2 to mock another process
//	MustNewTask(context.Background(),
//		"my_test_task_with_commit",
//		testDB,
//		testClient,
//		SetJobFunc(MyJobWithCommit),
//		SetParallelNum(10),
//		// SetHandleZombieJobsDstStatus(JobStatusTodo),// defualt: zombie
//		SetHandleZombieJobsPeriod(time.Second*5),
//		SetHandleZombieJobsTimeout(time.Hour*24),
//		SetHandleZombieJobsLimit(100),
//		SetTableName("task_test_table_with_commit"),
//		SetInitDB(false),
//		SetAutoCommit(false),
//	)
//
//	// run consumer 3 to mock another process
//	MustNewTask(context.Background(),
//		"my_test_task_with_commit",
//		testDB,
//		testClient,
//		SetJobFunc(MyJobWithCommit),
//		SetParallelNum(10),
//		// SetHandleZombieJobsDstStatus(JobStatusTodo), // defualt: zombie
//		SetHandleZombieJobsPeriod(time.Second*5),
//		SetHandleZombieJobsTimeout(time.Hour*24),
//		SetHandleZombieJobsLimit(100),
//		SetTableName("task_test_table_with_commit"),
//		SetInitDB(false),
//		SetAutoCommit(false),
//	)
//
//	// mock push data
//	totalJobs := 300
//	for i := 0; i < totalJobs; i++ {
//		commitTask.Push([]byte("hello"))
//	}
//
//	time.Sleep(time.Second * 30)
//
//	// check result
//	list := make([]*mysqlJob, 0)
//	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
//	defer cancel()
//	rows, err := testDB.QueryContext(ctx, fmt.Sprintf(queryJob, "task_test_table_with_commit"), "my_test_task_with_commit")
//
//	fmt.Println(err)
//
//	for rows.Next() {
//		job := new(mysqlJob)
//		// id, name, task_id, task_name, status, payload, failed_times, created_at, updated_at
//		err = rows.Scan(&job.Id, &job.Name, &job.TaskID, &job.TaskName, &job.Status, &job.Payload, &job.FailedTimes, &job.CommitID, &job.CreatedAt, &job.UpdatedAt)
//
//		list = append(list, job)
//	}
//	rows.Close()
//
//	failed := 0
//	todo := 0
//	done := 0
//	zombie := 0
//	for _, job := range list {
//		jobIDMap[job.Id] = job
//		switch JobStatus(job.Status) {
//		case JobStatusDone:
//			done++
//		case JobStatusFailed:
//			failed++
//		case JobStatusTodo:
//			todo++
//		case jobStatusZombie:
//			zombie++
//		}
//	}
//}
//
//func TestRePush(t *testing.T) {
//
//	var err error
//	commitTask1 := MustNewTask(context.Background(),
//		"my_test_task_with_commit",
//		testDB,
//		testClient,
//		SetParallelNum(0),
//		SetTableName("task_test_table_with_commit"),
//	)
//
//	// run consumer 2 to mock another process
//	commitTask2 := MustNewTask(context.Background(),
//		"my_test_task_with_commit",
//		testDB,
//		testClient,
//		SetParallelNum(0),
//		SetTableName("task_test_table_with_commit"),
//	)
//
//	successJob, failedJob, zombieJob := make([]*mysqlJob, 0), make([]*mysqlJob, 0), make([]*mysqlJob, 0)
//	for _, job := range jobIDMap {
//		if job.Status == "failed" {
//			failedJob = append(failedJob, job)
//		} else if job.Status == "done" {
//			successJob = append(successJob, job)
//		} else if job.Status == "zombie" && job.Id%99 != 0 { // 不希望repush id%99==0
//			zombieJob = append(zombieJob, job)
//		}
//	}
//
//	repushJobMap := make(map[int64]*mysqlJob)
//	fmt.Printf("successJob : %d", len(successJob))
//	fmt.Printf("failedJob : %d", len(failedJob))
//	fmt.Printf("zombieJob : %d", len(zombieJob))
//	for i := 0; i < 3; i++ {
//		err = commitTask2.RePush(successJob[i].Id)
//		if err != nil {
//			fmt.Printf("Repush error: %v", err)
//		}
//		repushJobMap[successJob[i].Id] = successJob[i]
//		fmt.Printf("Repush successJob ID: %d", successJob[i].Id)
//
//		err = commitTask2.RePush(failedJob[i].Id)
//		if err != nil {
//			fmt.Printf("Repush error: %v", err)
//		}
//		repushJobMap[failedJob[i].Id] = failedJob[i]
//		fmt.Printf("Repush failedJob ID: %d", failedJob[i].Id)
//
//		err = commitTask2.RePush(zombieJob[i].Id)
//		if err != nil {
//			fmt.Printf("Repush error: %v", err)
//		}
//		repushJobMap[zombieJob[i].Id] = zombieJob[i]
//		fmt.Printf("Repush zombieJob ID: %d", zombieJob[i].Id)
//	}
//	for i := 5; i < 10; i++ {
//		err = commitTask1.RePush(successJob[i].Id)
//		if err != nil {
//			fmt.Printf("Repush error: %v", err)
//		}
//		repushJobMap[successJob[i].Id] = successJob[i]
//		fmt.Printf("Repush successJob ID: %d", successJob[i].Id)
//
//		err = commitTask1.RePush(zombieJob[i].Id)
//		if err != nil {
//			fmt.Printf("Repush error: %v", err)
//		}
//		repushJobMap[zombieJob[i].Id] = zombieJob[i]
//		fmt.Printf("Repush zombieJob ID: %d", zombieJob[i].Id)
//	}
//
//	time.Sleep(time.Second * 30)
//
//	// check result
//	list := make([]*mysqlJob, 0)
//	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
//	defer cancel()
//	rows, err := testDB.QueryContext(ctx, fmt.Sprintf(queryJob, "task_test_table_with_commit"), "my_test_task_with_commit")
//
//	for rows.Next() {
//		job := new(mysqlJob)
//		// id, name, task_id, task_name, status, payload, failed_times, created_at, updated_at
//		err = rows.Scan(&job.Id, &job.Name, &job.TaskID, &job.TaskName, &job.Status, &job.Payload, &job.FailedTimes, &job.CommitID, &job.CreatedAt, &job.UpdatedAt)
//
//		list = append(list, job)
//	}
//	rows.Close()
//
//	failed := 0
//	todo := 0
//	done := 0
//	zombie := 0
//	for _, job := range list {
//		switch JobStatus(job.Status) {
//		case JobStatusDone:
//			done++
//		case JobStatusFailed:
//			failed++
//			if _, ok := repushJobMap[job.Id]; ok {
//			} else {
//			}
//		case JobStatusTodo:
//			todo++
//		case jobStatusZombie:
//			zombie++
//		}
//	}
//
//}
//

//
//// JobFunc example
//func MyJobWithCommit(ctx context.Context, id, fts, cID int64, payload []byte) error {
//
//	succ := false
//
//	defer func() {
//		if succ {
//			commitTask.Commit(id, cID, JobStatusDone)
//		} else {
//			if fts+1 >= 5 {
//				commitTask.Commit(id, cID, JobStatusFailed)
//			} else {
//				commitTask.Commit(id, cID, JobStatusTodo)
//			}
//		}
//	}()
//
//	fmt.Println("MyJob assigned job: ", id, string(payload))
//
//	// check for Idempotent
//	fmt.Println("check for idempotent...")
//
//	// then, do something
//	fmt.Println("do something ...")
//
//	// return if error with succ=false
//	if id%99 == 0 {
//		return fmt.Errorf("mock error")
//	}
//
//	succ = true
//
//	return nil
//
//}

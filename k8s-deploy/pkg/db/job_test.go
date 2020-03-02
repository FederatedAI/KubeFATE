package db

import (
	"context"
	"fate-cloud-agent/pkg/utils/logging"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

var jobJustAddedUuid string

func TestAddJob(t *testing.T) {
	InitConfigForTest()

	job := NewJob("cluster", "userid")
	JobUuid, error := Save(job)
	if error == nil {
		t.Log(JobUuid)
		jobJustAddedUuid = JobUuid
	}
}

func TestFindJobs(t *testing.T) {
	InitConfigForTest()
	job := &Job{}
	results, _ := Find(job)
	t.Log(ToJson(results))
}

func TestJobFindByUUID(t *testing.T) {
	InitConfigForTest()
	job := &Job{}
	results, _ := FindByUUID(job, jobJustAddedUuid)
	t.Log(ToJson(results))
}

func TestUpdateStatusByUUID(t *testing.T) {
	InitConfigForTest()
	t.Log("Update: " + jobJustAddedUuid)
	job := &Job{}
	result, error := FindByUUID(job, jobJustAddedUuid)
	if error == nil {
		job2Update := result.(Job)
		job2Update.Status = Success_j
		job2Update.EndTime = time.Now()
		UpdateByUUID(&job2Update, jobJustAddedUuid)
	}
	result, error = FindByUUID(job, jobJustAddedUuid)
	t.Log(ToJson(result))
}

func TestDeleteJobByUUID(t *testing.T) {
	InitConfigForTest()
	job := &Job{}
	DeleteByUUID(job, jobJustAddedUuid)
}

func TestFindJobList(t *testing.T) {
	InitConfigForTest()
	logging.InitLog()
	type args struct {
		args string
	}
	tests := []struct {
		name    string
		args    args
		want    []*Job
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "get job list",
			args:    args{args: ""},
			want:    make([]*Job, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JobFindList(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("JobFindList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i, v := range got {
				log.Info().Int("key", i).Interface("job", v).Msg("got")
			}
		})
	}
}

func TestFindJobByUUID(t *testing.T) {
	InitConfigForTest()
	logging.InitLog()
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "find db",
			args: args{
				uuid: "0c4da2a9-562b-4ce0-a564-46e318a85061",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JobFindByUUID(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("JobFindByUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("JobFindByUUID() = %+v", got)
		})
	}
}

func TestJobDeleteByUUID(t *testing.T) {
	InitConfigForTest()
	logging.InitLog()
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				uuid: "",
			},
			wantErr: true,
		},
		{
			name: "test",
			args: args{
				uuid: "2f75b214-6ea6-4cd6-9a52-14c4ceacb2b6",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := JobDeleteByUUID(tt.args.uuid); (err != nil) != tt.wantErr {
				t.Errorf("JobDeleteByUUID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJobStatus_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		s       JobStatus
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "",
			s:       Running_j,
			want:    []byte{34, 82, 117, 110, 110, 105, 110, 103, 34},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("JobStatus.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JobStatus.MarshalJSON() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestJobDeleteAll(t *testing.T) {
	InitConfigForTest()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	db, err := ConnectDb()
	if err != nil {
		log.Error().Err(err).Msg("ConnectDb")
	}
	collection := db.Collection(new(Job).getCollection())
	filter := bson.D{}
	r, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("DeleteMany")
	}
	if r.DeletedCount == 0 {
		log.Error().Msg("this record may not exist(DeletedCount==0)")
	}
	fmt.Println(r)
	return
}

func TestJob_timeOut(t *testing.T) {
	type fields struct {
		Uuid      string
		StartTime time.Time
		EndTime   time.Time
		Method    string
		Result    string
		ClusterId string
		Creator   string
		SubJobs   []string
		Status    JobStatus
		timeLimit time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			fields: fields{
				StartTime: time.Now(),
				EndTime:   time.Time{},
				timeLimit: 3 * time.Second,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := &Job{
				Uuid:      tt.fields.Uuid,
				StartTime: tt.fields.StartTime,
				EndTime:   tt.fields.EndTime,
				Method:    tt.fields.Method,
				Result:    tt.fields.Result,
				ClusterId: tt.fields.ClusterId,
				Creator:   tt.fields.Creator,
				SubJobs:   tt.fields.SubJobs,
				Status:    tt.fields.Status,
				timeLimit: tt.fields.timeLimit,
			}

			var i int

			for i < 10 {
				i++
				got := job.TimeOut()
				if got {
					t.Log("timeOut")
					return
				} else {
					t.Log("no timeOut")
				}
				time.Sleep(time.Second)
			}
			if i > 4 {
				t.Error("job.TimeOut() error")
			}
		})
	}
}

package mongo

import (
	"context"
	"fmt"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/pkg/errors"
)

const (
	jobCollectionName = "job"
)

type JobStore struct {
	Database *mongo.Database
}

type mongoJob struct {
	ObjectID      primitive.ObjectID  `bson:"_id,omitempty"`
	RepositoryID  primitive.ObjectID  `bson:"repository_id"`
	EnvironmentID *primitive.ObjectID `bson:"environment_id"`
	store.Job     `bson:",inline"`
}

type mongoJobUpdate struct {
	State     *shared.JobState `bson:"state,omitempty"`
	StoppedAt *time.Time       `bson:"stopped_at,omitempty"`
	StoppedBy *string          `bson:"stopped_by,omitempty"`
}

func (r *JobStore) Next() (*store.Job, error) {
	var job mongoJob
	err := r.
		Database.
		Collection(jobCollectionName).
		FindOneAndUpdate(
			context.Background(),
			bson.M{"state": shared.JobStateWaiting},
			bson.M{"$set": bson.M{"state": shared.JobStateProcessing, "started_at": time.Now()}},
		).
		Decode(&job)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, errors.Wrap(err, "error getting next available job")
		}
		return nil, nil
	}
	job.Job.ID = shared.JobID(job.ObjectID.Hex())
	job.Job.RepositoryID = store.RepositoryID(job.RepositoryID.Hex())
	if job.EnvironmentID != nil {
		e := shared.EnvironmentID(job.EnvironmentID.Hex())
		job.Job.EnvironmentID = &e
	}
	if job.EnvironmentID != nil {
		hex := shared.EnvironmentID(job.EnvironmentID.Hex())
		job.Job.EnvironmentID = &hex
	}

	return &job.Job, nil
}

func (r *JobStore) Get(id shared.JobID) (store.Job, error) {
	jobID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return store.Job{}, errors.Wrap(err, "error parsing id")
	}

	var mJob mongoJob
	err = r.
		Database.
		Collection(jobCollectionName).
		FindOne(
			context.Background(),
			bson.M{"_id": jobID},
		).Decode(&mJob)
	if err == mongo.ErrNoDocuments {
		return store.Job{}, store.ErrorNotFound
	}
	if err != nil {
		return store.Job{}, errors.Wrap(err, "error getting job")
	}

	mJob.Job.ID = shared.JobID(mJob.ObjectID.Hex())
	mJob.Job.RepositoryID = store.RepositoryID(mJob.RepositoryID.Hex())
	if mJob.EnvironmentID != nil {
		hex := shared.EnvironmentID(mJob.EnvironmentID.Hex())
		mJob.Job.EnvironmentID = &hex
	}
	return mJob.Job, nil
}

func (r *JobStore) Add(j store.Job) (shared.JobID, error) {
	mJob := mongoJob{Job: j}
	repoID, err := primitive.ObjectIDFromHex(string(j.RepositoryID))
	if err != nil {
		return shared.JobID(""), errors.Wrap(err, "error parsing id")
	}
	mJob.RepositoryID = repoID

	if j.EnvironmentID != nil {
		envID, err := primitive.ObjectIDFromHex(string(*j.EnvironmentID))
		if err != nil {
			return shared.JobID(""), errors.Wrap(err, "error parsing id")
		}
		mJob.EnvironmentID = &envID
	}

	result, err := r.
		Database.
		Collection(jobCollectionName).
		InsertOne(
			context.Background(),
			mJob,
		)
	if err != nil {
		return shared.JobID(""), errors.Wrap(err, "error adding job")
	}
	return shared.JobID(result.InsertedID.(primitive.ObjectID).String()), nil
}

func (r *JobStore) update(id shared.JobID, update mongoJobUpdate) error {
	objectID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return errors.Wrap(err, "error parsing id")
	}
	_, err = r.
		Database.
		Collection(jobCollectionName).
		UpdateOne(
			context.Background(),
			bson.M{"_id": objectID},
			bson.M{"$set": update},
		)
	return errors.Wrap(err, "error setting job state")
}

func (r *JobStore) UpdateStoppedAtByID(id shared.JobID, t time.Time) error {
	return r.update(id, mongoJobUpdate{StoppedAt: &t})
}

func (r *JobStore) UpdateStateByID(id shared.JobID, s shared.JobState) error {
	return r.update(id, mongoJobUpdate{State: &s})
}

func (r *JobStore) CancelByID(id shared.JobID, userID string) error {
	state := shared.JobStateCancelled
	return r.update(id, mongoJobUpdate{StoppedBy: &userID, State: &state})
}

func (r *JobStore) FilterByRepositoryID(
	repositoryID string,
	filter string,
	pageIndex int64,
	pageSize int64,
	sortColumn string,
	sortOrder int,
) (store.JobListPage, error) {
	page := store.JobListPage{}
	page.Jobs = []store.Job{}

	repositoryObjectID, err := primitive.ObjectIDFromHex(repositoryID)
	if err != nil {
		return page, errors.Wrap(err, "error parsing id")
	}

	bsonFilter := []bson.M{
		{"$match": bson.M{"repository_id": repositoryObjectID}},
		{"$match": bson.M{
			"$or": []bson.M{
				{"commit.branch": bsonx.Regex(fmt.Sprintf(".*%s.*", filter), "")},
				{"commit.revision": bsonx.Regex(fmt.Sprintf(".*%s.*", filter), "")},
				{"started_by": bsonx.Regex(fmt.Sprintf(".*%s.*", filter), "")},
			},
		},
		},
	}

	decoder, err := r.
		Database.
		Collection(jobCollectionName).Aggregate(
		context.Background(),
		append(
			bsonFilter,
			bson.M{"$sort": bson.M{string(sortColumn): sortOrder}},
			bson.M{"$skip": pageIndex * pageSize},
			bson.M{"$limit": pageSize},
		),
	)
	if err != nil {
		return page, errors.Wrap(err, "error fetching jobs")
	}

	for decoder.Next(context.Background()) {
		var r mongoJob
		err = decoder.Decode(&r)
		if err != nil {
			return page, errors.Wrap(err, "error decoding job")
		}
		r.Job.ID = shared.JobID(r.ObjectID.Hex())
		r.Job.RepositoryID = store.RepositoryID(r.RepositoryID.Hex())
		if r.Job.EnvironmentID != nil {
			hex := shared.EnvironmentID(r.EnvironmentID.Hex())
			r.Job.EnvironmentID = &hex
		}

		page.Jobs = append(page.Jobs, r.Job)
	}

	decoder, err = r.
		Database.
		Collection(jobCollectionName).Aggregate(
		context.Background(),
		append(
			bsonFilter,
			bson.M{"$count": "job_count"},
		),
	)
	if err != nil {
		return page, errors.Wrap(err, "error counting jobs")
	}

	if decoder.Next(context.Background()) {
		err = decoder.Decode(&page)
		if err != nil {
			return page, errors.Wrap(err, "error decoding job count")
		}
	}
	return page, nil
}

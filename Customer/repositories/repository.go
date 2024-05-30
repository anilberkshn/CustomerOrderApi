package repositories

import (
	"CustomerOrderApi/Customer/entities"
	"context"
	"fmt"
	"github.com/erenkaratas99/COApiCore/pkg/customErrors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository() *Repository {

	databaseURL := "mongodb://127.0.0.1:27017"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(databaseURL))
	if err != nil {
		log.Fatal("Hata : " + err.Error())
	}

	//clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	//client, err := mongo.Connect(context.Background(), clientOptions)
	//if err != nil {log.Fatal(err)}

	collection := client.Database("goCustomers").Collection("goCustomers")

	return &Repository{collection}
}

func (r *Repository) InsertCustomer(customerReq *entities.CustomerRequestModel) (*string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	id := uuid.New()
	customerId := id.String()
	customerReq.Address.CustomerId = customerId
	timeNow := time.Now().Format(time.RFC3339)
	c := bson.M{
		"_id":        customerId,
		"first_name": customerReq.FirstName,
		"last_name":  customerReq.LastName,
		"email":      customerReq.Email,
		"phone":      customerReq.Phone,
		"address":    customerReq.Address,
		"created_at": timeNow,
		"updated_at": timeNow,
	}
	res, err := r.collection.InsertOne(ctx, c)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	insertedId := res.InsertedID.(string)
	return &insertedId, nil
}

func (r *Repository) GetAllCustomers(l, o int64, getTotalCount bool) ([]*entities.CustomerResponseModel, *int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	//todo : tekrar

	filter := bson.M{}
	opts := options.Find().SetLimit(l).SetSkip(o)
	var err error
	var customer []*entities.CustomerResponseModel

	var wg sync.WaitGroup
	var totalCountG int
	wg.Add(1)
	go func() {
		defer wg.Done()
		cur, err := r.collection.Find(ctx, filter, opts)
		if err != nil {
			return
		}
		defer cur.Close(ctx)
		err = cur.All(ctx, &customer)
	}()
	if getTotalCount {
		totalCount64, _ := r.collection.CountDocuments(ctx, filter)
		totalCountG = int(totalCount64)
	}
	wg.Wait()
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	return customer, &totalCountG, nil
}

func (r *Repository) CustomerGetById(id string, getTotalCount bool) (*entities.CustomerResponseModel, *int, error) {
	filter := bson.M{"_id": id}
	customer := entities.CustomerResponseModel{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := r.collection.FindOne(ctx, filter).Decode(&customer)

	if err != nil {
		fmt.Println(err)
		return nil, nil, customErrors.DocNotFound
	}
	var totalCountG int
	if getTotalCount {
		totalCount64, _ := r.collection.CountDocuments(ctx, filter)
		totalCountG = int(totalCount64)
	}
	return &customer, &totalCountG, nil
}

func (r *Repository) Save(data interface{}) {
	// MongoDB'ye veri kaydet
}

package customer

import (
	"fmt"
	"goworkshop2/storage"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MongoCustomerTester struct {
	Service  MongoCustomerService
	Database storage.MongoDb
}

const EMAIL_FOR_TEST = "A@a.com"
const PASSWORD_FOR_TEST = "1234"
const NEWPASSWORD_FOR_TEST = "4567"
const NEWNAME_FOR_TEST = "New world"
const NAME_FOR_TEST = "Hello world"

func setup() MongoCustomerTester {
	fmt.Println("SETUP==============================")
	os.Setenv("MONGODB_USERNAME", "user2")
	os.Setenv("MONGODB_PASSWORD", "2222")
	os.Setenv("MONGODB_DB_NAME", "goworkshop2")
	os.Setenv("MONGODB_ENDPOINT", "localhost")

	dbConfig := storage.MongoDbConfiguration{
		ConnectionTimeout:        10,
		ConnectionStringTemplate: storage.DEFAULT_CONNECTION_STRING_FORMAT,
	}
	db, _ := storage.CreateDatabase(&dbConfig)

	service := MongoCustomerService{
		Resource: db.Resource,
	}

	err := db.DropAll()
	if err != nil {
		fmt.Println("setup but got error on drop :", err)
	}

	return MongoCustomerTester{
		Service:  service,
		Database: db,
	}
}

func teardown(tester MongoCustomerTester) {
	tester.Database.Close()
}

func Test_Register_firsttime_should_success(t *testing.T) {
	env := setup()
	defer teardown(env)

	err := env.Service.Register("A@a.com", "1234", "hi sir")

	assert.Nil(t, err)
}

func Test_Register_exist_user_should_fail(t *testing.T) {
	env := setup()
	defer teardown(env)

	env.Service.Register("A@a.com", "1234", "hi sir")
	err := env.Service.Register("A@a.com", "1234", "hi sir")

	assert.NotNil(t, err)
}

func Test_Login_WithNoneExistUser_should_be_fail(t *testing.T) {
	env := setup()
	defer teardown(env)

	result, _ := env.Service.Login("A@a.com", "1234")

	assert.Equal(t, false, result)
}

func Test_Login_WithWrongPassword_should_be_fail(t *testing.T) {
	env := setup()
	defer teardown(env)

	env.Service.Register(EMAIL_FOR_TEST, PASSWORD_FOR_TEST, NAME_FOR_TEST)
	result, _ := env.Service.Login(EMAIL_FOR_TEST, "PASSWORD_FOR_TEST")

	assert.Equal(t, false, result)
}

func Test_Login_should_be_success(t *testing.T) {
	env := setup()
	defer teardown(env)

	env.Service.Register(EMAIL_FOR_TEST, PASSWORD_FOR_TEST, NAME_FOR_TEST)
	result, _ := env.Service.Login(EMAIL_FOR_TEST, PASSWORD_FOR_TEST)

	assert.Equal(t, true, result)
}

func Test_ChangePassword_WithNonExistUser_should_fail(t *testing.T) {
	env := setup()
	defer teardown(env)

	result := env.Service.ChangePassword(EMAIL_FOR_TEST, PASSWORD_FOR_TEST, NEWPASSWORD_FOR_TEST)

	assert.NotNil(t, result)
}

func Test_ChangePassword_WithOldPasswordNotCorrect_should_fail(t *testing.T) {
	env := setup()
	defer teardown(env)

	env.Service.Register(EMAIL_FOR_TEST, PASSWORD_FOR_TEST, NAME_FOR_TEST)
	err := env.Service.ChangePassword(EMAIL_FOR_TEST, "PASSWORD_FOR_TEST", NEWPASSWORD_FOR_TEST)

	assert.NotNil(t, err)
}

func Test_ChangePassword_should_success(t *testing.T) {
	env := setup()
	defer teardown(env)

	env.Service.Register(EMAIL_FOR_TEST, PASSWORD_FOR_TEST, NAME_FOR_TEST)
	err := env.Service.ChangePassword(EMAIL_FOR_TEST, PASSWORD_FOR_TEST, NEWPASSWORD_FOR_TEST)

	assert.Nil(t, err)
}

func Test_EditProfile_should_success(t *testing.T) {
	env := setup()
	defer teardown(env)

	env.Service.Register(EMAIL_FOR_TEST, PASSWORD_FOR_TEST, NAME_FOR_TEST)
	err := env.Service.UpdateProfile(EMAIL_FOR_TEST, NEWNAME_FOR_TEST)

	assert.Nil(t, err)
}

func Test_GetProfile_should_success(t *testing.T) {
	env := setup()
	defer teardown(env)

	env.Service.Register(EMAIL_FOR_TEST, PASSWORD_FOR_TEST, NAME_FOR_TEST)
	result, _ := env.Service.GetProfile(EMAIL_FOR_TEST)

	assert.Equal(t, NAME_FOR_TEST, result.Name)
}

func Test_GetProfile_OfNonExistUser_should_fail(t *testing.T) {
	env := setup()
	defer teardown(env)

	result, err := env.Service.GetProfile(EMAIL_FOR_TEST)

	assert.NotNil(t, err)
	assert.Equal(t, result, User{})
}

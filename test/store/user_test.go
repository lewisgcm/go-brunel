package store

import (
	"go-brunel/internal/pkg/server/security"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/test"
	"testing"
	"time"
)

func TestAddUser(t *testing.T) {
	suite := setup(t)

	for _, userStore := range suite.userStores {
		now := time.Now()
		user, err := userStore.AddOrUpdate(store.User{
			Username:  "username",
			Email:     "email",
			Name:      "name",
			AvatarURL: "avatar",
			Role:      security.UserRoleAdmin,
		})

		if err != nil {
			t.Fatalf("could not create user: %s", err)
		}

		if e := userStore.Delete(user.Username, true); e != nil {
			t.Fatalf("could not delete test user: %s", e)
		}

		test.ExpectString(t, "username", user.Username)
		test.ExpectString(t, "email", user.Email)
		test.ExpectString(t, "name", user.Name)
		test.ExpectString(t, "avatar", user.AvatarURL)
		test.ExpectString(t, string(security.UserRoleAdmin), string(user.Role))

		if now.Unix() < user.CreatedAt.Unix() {
			t.Errorf("created at date is invalid: %s < %s", user.CreatedAt.String(), now.String())
		}
	}
}

func TestAddUserDefaultRole(t *testing.T) {
	suite := setup(t)

	for _, userStore := range suite.userStores {
		user, err := userStore.AddOrUpdate(store.User{
			Username:  "username",
			Email:     "email",
			Name:      "name",
			AvatarURL: "avatar",
		})

		if err != nil {
			t.Fatalf("could not create user: %s", err)
		}

		if e := userStore.Delete(user.Username, true); e != nil {
			t.Fatalf("could not delete test user: %s", e)
		}

		test.ExpectString(t, string(security.UserRoleReader), string(user.Role))
	}
}

func TestUpdateUser(t *testing.T) {
	suite := setup(t)

	for _, userStore := range suite.userStores {
		user, err := userStore.AddOrUpdate(store.User{
			Username:  "username",
			Email:     "email",
			Name:      "name",
			AvatarURL: "avatar",
			Role:      security.UserRoleAdmin,
		})

		if err != nil {
			t.Fatalf("could not create user: %s", err)
		}

		updatedUser, err := userStore.AddOrUpdate(store.User{
			Username:  "username",
			Email:     "email2",
			Name:      "name2",
			AvatarURL: "avatar2",
			Role:      security.UserRoleReader,
		})

		if e := userStore.Delete(user.Username, true); e != nil {
			t.Fatalf("could not delete test user: %s", e)
		}

		if err != nil {
			t.Fatalf("could not update user: %s", err)
		}

		test.ExpectString(t, "username", updatedUser.Username)
		test.ExpectString(t, "email2", updatedUser.Email)
		test.ExpectString(t, "name2", updatedUser.Name)
		test.ExpectString(t, "avatar2", updatedUser.AvatarURL)
		test.ExpectString(t, string(security.UserRoleReader), string(updatedUser.Role))

		if !user.CreatedAt.Equal(updatedUser.CreatedAt) {
			t.Errorf("created at date is invalid: %s == %s", user.CreatedAt.String(), updatedUser.CreatedAt.String())
		}
	}
}

func TestFilterUser(t *testing.T) {
	suite := setup(t)

	for _, userStore := range suite.userStores {
		user, err := userStore.AddOrUpdate(store.User{
			Username:  "username",
			Email:     "email",
			Name:      "name",
			AvatarURL: "avatar",
			Role:      security.UserRoleAdmin,
		})

		if err != nil {
			t.Fatalf("could not create user: %s", err)
		}

		users, err := userStore.Filter("")
		if err != nil || len(users) == 0 {
			t.Fatalf("invalid results returned for search '%s'", "")
		}

		users, err = userStore.Filter("UsERNAMe")
		if err != nil || len(users) == 0 {
			t.Fatalf("invalid results returned for search '%s'", "UsERNAMe")
		}

		users, err = userStore.Filter("EmAIl")
		if err != nil || len(users) == 0 {
			t.Fatalf("invalid results returned for search '%s'", "EmAIl")
		}

		users, err = userStore.Filter("ddddddd")
		if err != nil || len(users) != 0 {
			t.Fatalf("invalid results returned for search '%s'", "ddddddd")
		}

		if e := userStore.Delete(user.Username, true); e != nil {
			t.Fatalf("could not delete test user: %s", e)
		}
	}
}

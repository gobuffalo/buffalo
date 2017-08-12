package pop_test

import (
	"testing"

	"github.com/markbates/pop"
	"github.com/markbates/pop/nulls"
	"github.com/stretchr/testify/require"
)

func Test_ValidateAndSave(t *testing.T) {
	r := require.New(t)
	validationLogs = []string{}
	transaction(func(tx *pop.Connection) {
		car := &ValidatableCar{Name: "VW"}
		verrs, err := tx.ValidateAndSave(car)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 2)
		r.Equal([]string{"Validate", "ValidateSave"}, validationLogs)
		r.NotZero(car.ID)
		r.NotZero(car.CreatedAt)

		validationLogs = []string{}
		car = &ValidatableCar{Name: ""}
		verrs, err = tx.ValidateAndSave(car)
		r.NoError(err)
		r.True(verrs.HasAny())
		r.Len(validationLogs, 2)
		errs := verrs.Get("name")
		r.Len(errs, 1)

		validationLogs = []string{}
		ncar := &NotValidatableCar{Name: ""}
		verrs, err = tx.ValidateAndSave(ncar)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 0)
	})
}

func Test_ValidateAndCreate(t *testing.T) {
	r := require.New(t)
	validationLogs = []string{}
	transaction(func(tx *pop.Connection) {
		car := &ValidatableCar{Name: "VW"}
		verrs, err := tx.ValidateAndCreate(car)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 2)
		r.Equal([]string{"Validate", "ValidateCreate"}, validationLogs)
		r.NotZero(car.ID)
		r.NotZero(car.CreatedAt)

		validationLogs = []string{}
		car = &ValidatableCar{Name: ""}
		verrs, err = tx.ValidateAndSave(car)
		r.NoError(err)
		r.True(verrs.HasAny())
		r.Len(validationLogs, 2)
		errs := verrs.Get("name")
		r.Len(errs, 1)

		validationLogs = []string{}
		ncar := &NotValidatableCar{Name: ""}
		verrs, err = tx.ValidateAndCreate(ncar)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 0)
	})
}

func Test_ValidateAndUpdate(t *testing.T) {
	r := require.New(t)
	validationLogs = []string{}
	transaction(func(tx *pop.Connection) {
		car := &ValidatableCar{Name: "VW"}
		verrs, err := tx.ValidateAndCreate(car)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 2)
		r.Equal([]string{"Validate", "ValidateCreate"}, validationLogs)
		r.NotZero(car.ID)
		r.NotZero(car.CreatedAt)

		validationLogs = []string{}
		car.Name = ""
		verrs, err = tx.ValidateAndUpdate(car)
		r.NoError(err)
		r.True(verrs.HasAny())
		r.Len(validationLogs, 2)
		errs := verrs.Get("name")
		r.Len(errs, 1)

		validationLogs = []string{}
		ncar := &NotValidatableCar{Name: ""}
		verrs, err = tx.ValidateAndCreate(ncar)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 0)

		validationLogs = []string{}
		ncar.Name = ""
		verrs, err = tx.ValidateAndUpdate(ncar)
		r.NoError(err)
		r.False(verrs.HasAny())
		r.Len(validationLogs, 0)
	})
}

func Test_Exec(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		user := User{Name: nulls.NewString("Mark 'Awesome' Bates")}
		tx.Create(&user)

		ctx, _ := tx.Count(user)
		a.Equal(1, ctx)

		q := tx.RawQuery("delete from users where id = ?", user.ID)
		err := q.Exec()
		a.NoError(err)

		ctx, _ = tx.Count(user)
		a.Equal(0, ctx)
	})
}

func Test_Save(t *testing.T) {
	r := require.New(t)
	transaction(func(tx *pop.Connection) {
		u := &User{Name: nulls.NewString("Mark")}
		r.Zero(u.ID)
		tx.Save(u)
		r.NotZero(u.ID)

		uat := u.UpdatedAt.UnixNano()

		tx.Save(u)
		r.NotEqual(uat, u.UpdatedAt.UnixNano())
	})
}

func Test_Create(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		count, _ := tx.Count(&User{})
		user := User{Name: nulls.NewString("Mark 'Awesome' Bates")}
		err := tx.Create(&user)
		a.NoError(err)
		a.NotEqual(user.ID, 0)

		ctx, _ := tx.Count(&User{})
		a.Equal(count+1, ctx)

		u := User{}
		q := tx.Where("name = ?", "Mark 'Awesome' Bates")
		err = q.First(&u)
		a.NoError(err)
		a.Equal(user.Name.String, "Mark 'Awesome' Bates")
	})
}

func Test_Create_UUID(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		count, _ := tx.Count(&Song{})
		song := Song{Title: "Automatic Buffalo"}
		err := tx.Create(&song)
		a.NoError(err)
		a.NotZero(song.ID)

		ctx, _ := tx.Count(&Song{})
		a.Equal(count+1, ctx)

		u := Song{}
		q := tx.Where("title = ?", "Automatic Buffalo")
		err = q.First(&u)
		a.NoError(err)
	})
}

func Test_Create_Timestamps(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		user := User{Name: nulls.NewString("Mark 'Awesome' Bates")}
		a.Zero(user.CreatedAt)
		a.Zero(user.UpdatedAt)

		err := tx.Create(&user)
		a.NoError(err)

		a.NotZero(user.CreatedAt)
		a.NotZero(user.UpdatedAt)

		friend := Friend{FirstName: "Ross", LastName: "Gellar"}
		err = tx.Create(&friend)
		a.NoError(err)
	})
}

func Test_Update(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		user := User{Name: nulls.NewString("Mark")}
		tx.Create(&user)

		a.NotZero(user.CreatedAt)
		a.NotZero(user.UpdatedAt)

		user.Name.String = "Marky"
		err := tx.Update(&user)
		a.NoError(err)

		tx.Reload(&user)
		a.Equal(user.Name.String, "Marky")
	})
}

func Test_Update_UUID(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		r := require.New(t)

		song := Song{Title: "Automatic Buffalo"}
		err := tx.Create(&song)
		r.NoError(err)

		r.NotZero(song.CreatedAt)
		r.NotZero(song.UpdatedAt)

		song.Title = "Hum"
		err = tx.Update(&song)
		r.NoError(err)

		err = tx.Reload(&song)
		r.NoError(err)
		r.Equal("Hum", song.Title)
	})
}

func Test_Destroy(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		a := require.New(t)

		count, err := tx.Count("users")
		user := User{Name: nulls.NewString("Mark")}
		err = tx.Create(&user)
		a.NoError(err)
		a.NotEqual(user.ID, 0)

		ctx, err := tx.Count("users")
		a.Equal(count+1, ctx)

		err = tx.Destroy(&user)
		a.NoError(err)

		ctx, _ = tx.Count("users")
		a.Equal(count, ctx)
	})
}

func Test_Destroy_UUID(t *testing.T) {
	transaction(func(tx *pop.Connection) {
		r := require.New(t)

		count, err := tx.Count("songs")
		song := Song{Title: "Automatic Buffalo"}
		err = tx.Create(&song)
		r.NoError(err)
		r.NotZero(song.ID)

		ctx, err := tx.Count("songs")
		r.Equal(count+1, ctx)

		err = tx.Destroy(&song)
		r.NoError(err)

		ctx, _ = tx.Count("songs")
		r.Equal(count, ctx)
	})
}

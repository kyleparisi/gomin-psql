package app

type AppUser struct {
  Id int64
  Name string
  Email string
  Password string
}

type DjangoMigrations struct {
  App string
  Name string
}

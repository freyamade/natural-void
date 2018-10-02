# natural-void
A podcast site made specifically for the ExVo D&amp;D games.

That being said, I'll be making it as easy as possible for others to use so feel free :)

Made with Go, Bulma, TypeScript and &lt;3

## Environment Variables
The following environment variables are used in the configuration of the system;

### Postgres
- `POSTGRES_HOST`
    - The host name of the postgres server
    - Defaults to `localhost`
- `POSTGRES_PORT`
    - The port of the postgres server
    - Defaults to `5432`
- `POSTGRES_USER`
    - The username used to connect to postgres
    - Defaults to `postgres`
- `POSTGRES_PASS`
    - The password used to authenticate the postgres user
    - Defaults to `postgres`
- `POSTGRES_NAME`
    - This name of the database being used for the system
    - Defaults to `postgres`
- `POSTGRES_SECURE`
    - Flag stating whether or not the connection to postgres should be done with SSL
    - Defaults to `false`

### LDAP
- `LDAP_HOST`
    - The host name of the ldap server
    - Defaults to `localhost`
- `LDAP_PORT`
    - The port of the ldap server
    - Defaults to `389`
- `LDAP_USER`
    - The username used to connect to ldap
    - Defaults to `cn=admin,dc=naturalvoid,dc=com`
- `LDAP_PASS`
    - The password used to authenticate the ldap user
    - Defaults to `nv`
- `LDAP_SECURE`
    - Flag stating whether or not the connection to ldap should be done with SSL
    - Defaults to `false`

### Other
- `SECRET_KEY`
    - This is used to generate a secret key for secrets
    - Defaults to `replacethiswithanactualsecretkey`
- `PRODUCTION`
    - A flag stating whether or not the server is being run in a production environment
    - Defaults to `false`

## Development Mode
To develop this project further you'll need to run an LDAP and Postgres server locally.

A `docker-compose.yml` file has been supplied to set these up, and by default the conf is set to use these.

Run `docker-compose up -d` and also `go run seed.go` to pre-populate the test DB with some records.

A user has been set up for the ldap server with the following details;

- Username: `crnbrdrck`
- Password: `password`

## Custom Colouring
It is completely possible to customize the colour scheme for the website by adding a very basic scss file.

In order to colour the website, I come up with a 2 colour scheme.
This is expanded to a 5 colour scheme using websites like ColorMind, Coolors, etc.

The 5 colours are then used as `light-shade`, `light-accent`, `main-brand`, `dark-accent`, and `dark-shade`.
The 2 colours are also set as `scheme-primary` and `scheme-secondary` depending on which is your main and which is your secondary colour.

If you look at `scss/dread.scss` you'll see an example of a colour scheme;

```scss
// Set the 5 main colours
$light-shade: #B9ADA9;
$light-accent: #656993;
$main-brand: #2C0047;
$dark-accent: #600;
$dark-shade: #1A141D;

// Set the notification colours
$primary: #2c0047;
$primary-invert: $light-shade;
$info: #19141d;
$info-invert: white;
$success: #427a4d;
$success-invert: white;
$warning: #c06a15;
$warning-invert: white;
$danger: #f44336;
$danger-invert: white;

// Set the primary and secondary colours (they don't have to be from the ones defined above)
$scheme-primary: $main-brand;
$scheme-primary-invert: $light-shade;
$scheme-secondary: $dark-accent;
$scheme-secondary-invert: $light-shade;

@import "./_main.scss"
```

All you have to do is create an scss file similar to this, name it whatever you want and save it in the `scss` folder.
Then by running `npm run css-build`, you will get a file in `static/css` with the same name as the `scss` file you made.

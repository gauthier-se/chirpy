{ pkgs, lib, config, inputs, ... }:

{
  languages.go.enable = true;

  # .env is loaded by direnv (.envrc), whose parser strips quotes correctly;
  # devenv's own dotenv keeps them literal. Just silence the suggestion.
  dotenv.disableHint = true;

  services.postgres = {
    enable = true;
    listen_addresses = "127.0.0.1";
    package = pkgs.postgresql_15;
    initialDatabases = [
      { name = "postgres"; }
    ];
  };

  enterShell = ''
    echo "Go & Postgres environment ready!"
    echo "Run 'devenv up' in a separate terminal to start the database."
  '';
}

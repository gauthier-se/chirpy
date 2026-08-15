{ pkgs, lib, config, inputs, ... }:

{
  languages.go.enable = true;

  services.postgres = {
    enable = true;
    listen_addresses = "127.0.0.1";
    package = pkgs.postgresql_15;
    initialDatabases = [
      { name = "postgres"; }
    ];
  };

  enterShell = ''
    echo "Environnement Go & Postgres prêt !"
    echo "Lance 'devenv up' dans un terminal séparé pour démarrer la base de données."
  '';
}

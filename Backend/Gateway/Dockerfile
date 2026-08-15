FROM kong:3.7

# db-less: config comes from KONG_DECLARATIVE_CONFIG (see docker-compose.yml).
#
# When custom-auth/ has a real Lua plugin, bake it in and enable it:
#   USER root
#   COPY custom-auth /usr/local/share/lua/5.1/kong/plugins/custom-auth
#   USER kong
# then set KONG_PLUGINS=bundled,custom-auth in docker-compose.yml

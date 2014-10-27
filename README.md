Run

    cp ./run.sh.example ./run.sh

and edit the keys. Then run

    ./run.sh

Failed approaches can be found in `1` and `2`.

Please note that this code implements a minimal OAuth2 exchange for "installed
applications". E.g. it does not use any Google API client libraries, and does
not store the refresh token.

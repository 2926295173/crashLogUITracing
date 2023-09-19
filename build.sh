(
    cd frontend
    pnpm build
)

cp -r frontent/dist/* backend/static/

(
    cd backend
    sh build.sh
)
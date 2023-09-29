(
    cd frontend
    pnpm build
)

cp -r frontend/dist/* backend/static/

(
    cd backend
    sh build.sh
)
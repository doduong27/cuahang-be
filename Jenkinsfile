pipeline {
    agent any

    parameters {
        choice(
            name: 'TARGET_ENV',
            choices: ['dev', 'stag', 'prod'],
            description: '🎯 Chọn môi trường bạn muốn triển khai Backend (Go):'
        )
        choice(
            name: 'ACTION',
            choices: ['Build & Deploy', 'Deploy Only (Promote)'],
            description: '⚙️ Chọn hành động (Build & Deploy: build code mới; Deploy Only: đưa tag đã có sang môi trường này):'
        )
        string(
            name: 'IMAGE_TAG',
            defaultValue: 'v1.0.0',
            description: '🏷️ Tag phiên bản SemVer (ví dụ: v1.0.0, v1.1.0, v2.0.0):'
        )
    }

    environment {
        APP_NAME        = 'cuahang-be'
        DOCKER_HUB_USER = 'dvduong'
        IMAGE_NAME      = "${DOCKER_HUB_USER}/${APP_NAME}"
        GITOPS_REPO_SSH = 'git@github.com:doduong27/cuahang-gitops.git'
    }

    stages {
        stage('1. Khởi tạo & Xác định phiên bản') {
            steps {
                script {
                    if (params.ACTION == 'Build & Deploy') {
                        env.FINAL_TAG = params.IMAGE_TAG ?: "build-${BUILD_NUMBER}-${GIT_COMMIT.take(7)}"
                    } else {
                        env.FINAL_TAG = params.IMAGE_TAG
                    }
                    echo "=========================================================="
                    echo "🚀 KẾ HOẠCH TRIỂN KHAI BACKEND (Go Gin + PostgreSQL):"
                    echo "   - Ứng dụng: ${APP_NAME}"
                    echo "   - Môi trường đích: [ ${params.TARGET_ENV.toUpperCase()} ]"
                    echo "   - Hành động: ${params.ACTION}"
                    echo "   - Phiên bản Image Tag: ${env.FINAL_TAG}"
                    echo "=========================================================="
                }
            }
        }

        stage('2. Quét bảo mật & Unit Test (Chỉ khi Build)') {
            when {
                expression { return params.ACTION == 'Build & Deploy' }
            }
            steps {
                echo '🛡️ 1. Quét rò rỉ secret trong mã nguồn (Gitleaks)...'
                sh '''
                    if command -v gitleaks &> /dev/null; then
                        gitleaks detect --source . --verbose --no-git || true
                    else
                        echo "Gitleaks chưa cài, bỏ qua quét secret"
                    fi
                '''

                echo '🛡️ 2. Phân tích tĩnh mã nguồn Go bằng Semgrep (SAST)...'
                sh '''
                    if command -v semgrep &> /dev/null; then
                        semgrep --config p/golang --error . || true
                    else
                        echo "Semgrep chưa cài, bỏ qua SAST"
                    fi
                '''

                echo '🧪 3. Chạy Unit Test Go...'
                sh '''
                    export GOPROXY=https://proxy.golang.org,https://goproxy.io,direct
                    export GOSUMDB=off
                    go test -v ./... || true
                '''
            }
        }

        stage('3. Build & Push Docker Image') {
            when {
                expression { return params.ACTION == 'Build & Deploy' }
            }
            steps {
                echo "🐳 Đang Build Docker Image ${IMAGE_NAME}:${env.FINAL_TAG}..."
                sh "docker build -t ${IMAGE_NAME}:${env.FINAL_TAG} -t ${IMAGE_NAME}:latest -t ${IMAGE_NAME}:${params.TARGET_ENV}-latest ."

                echo '🛡️ Quét lỗ hổng bảo mật Docker Image bằng Trivy...'
                sh """
                    if command -v trivy &> /dev/null; then
                        trivy image --severity CRITICAL --exit-code 0 ${IMAGE_NAME}:${env.FINAL_TAG} || true
                    else
                        echo "Trivy chưa cài, bỏ qua quét Image"
                    fi
                """

                echo "🚀 Đẩy image lên Docker Hub..."
                withCredentials([usernamePassword(credentialsId: 'docker-hub-credentials', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                    sh """
                        echo "\$DOCKER_PASS" | docker login -u "\$DOCKER_USER" --password-stdin
                        docker push ${IMAGE_NAME}:${env.FINAL_TAG}
                        docker push ${IMAGE_NAME}:latest
                        docker push ${IMAGE_NAME}:${params.TARGET_ENV}-latest
                    """
                }
            }
        }

        stage('4. Cập nhật GitOps Manifest') {
            steps {
                echo "🔄 Đang cập nhật Kustomize Manifest cho môi trường: ${params.TARGET_ENV}..."
                sh """
                    rm -rf gitops-tmp
                    GIT_SSH_COMMAND="ssh -o StrictHostKeyChecking=no" git clone ${GITOPS_REPO_SSH} gitops-tmp
                    cd gitops-tmp/projects/cuahang/overlays/${params.TARGET_ENV}
                    
                    # Sử dụng Kustomize CLI để cập nhật image tag backend chuẩn xác
                    kustomize edit set image cuahang-be=docker.io/${IMAGE_NAME}:${env.FINAL_TAG}
                    
                    git config user.name "jenkins-bot"
                    git config user.email "jenkins@lan"
                    
                    if git status --porcelain | grep -q 'kustomization.yaml'; then
                        git add kustomization.yaml
                        git commit -m "ci(${params.TARGET_ENV}): promote ${APP_NAME} to ${env.FINAL_TAG} via Jenkins [skip ci]"
                        GIT_SSH_COMMAND="ssh -o StrictHostKeyChecking=no" git push origin main
                        echo "✅ Đã commit và push cập nhật GitOps Backend thành công!"
                    else
                        echo "ℹ️ Tag ${env.FINAL_TAG} đã tồn tại trong kustomization.yaml của ${params.TARGET_ENV}, không có thay đổi mới."
                    fi
                """
            }
        }

        stage('5. Kiểm tra trạng thái triển khai') {
            steps {
                script {
                    echo "🎉 Đã hoàn tất đẩy cấu hình Backend vào GitOps!"
                    echo "Argo CD sẽ tự động đồng bộ (hoặc bấm Sync trên Argo CD nếu là Production)."
                    def targetHost = "${params.TARGET_ENV}.cuahang.lan"
                    if (params.TARGET_ENV == 'prod') {
                        targetHost = "cuahang.lan"
                    }
                    echo "🌐 Endpoint kiểm tra API Health: http://${targetHost}/healthz"
                    echo "🌐 Endpoint kiểm tra API Products: http://${targetHost}/api/products"
                }
            }
        }
    }

    post {
        success {
            echo "✅ Pipeline Backend hoàn thành xuất sắc! Môi trường [ ${params.TARGET_ENV.toUpperCase()} ] đã được đồng bộ với Tag [ ${env.FINAL_TAG} ]."
        }
        failure {
            echo "❌ Pipeline Backend gặp lỗi! Vui lòng kiểm tra log chi tiết."
        }
    }
}

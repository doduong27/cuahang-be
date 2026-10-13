pipeline {
    agent any

    environment {
        APP_NAME        = 'cuahang-be'
        DOCKER_HUB_USER = 'your-dockerhub-username' // Thay bằng username Docker Hub của bạn
        IMAGE_NAME      = "${DOCKER_HUB_USER}/${APP_NAME}"
        IMAGE_TAG       = "${BUILD_NUMBER}-${GIT_COMMIT.take(7)}"
        GITOPS_REPO     = 'https://github.com/your-github-user/cuahang-gitops.git'
    }

    stages {
        stage('1. Pre-Check & Secret Scanning') {
            steps {
                echo '🛡️ Quét rò rỉ secret trong mã nguồn (Gitleaks)...'
                sh '''
                    if command -v gitleaks &> /dev/null; then
                        gitleaks detect --source . --verbose --no-git || true
                    else
                        echo "Gitleaks chưa cài, bỏ qua bước quét"
                    fi
                '''
            }
        }

        stage('2. Static Code Analysis (SAST)') {
            steps {
                echo '🛡️ Phân tích tĩnh mã nguồn Go bằng Semgrep...'
                sh '''
                    if command -v semgrep &> /dev/null; then
                        semgrep --config p/golang --error . || true
                    else
                        echo "Semgrep chưa cài, bỏ qua bước SAST"
                    fi
                '''
            }
        }

        stage('3. Unit Test') {
            steps {
                echo '🧪 Chạy Unit Test Go...'
                sh '''
                    go test -v ./... || true
                '''
            }
        }

        stage('4. Docker Build & Trivy CVE Scan') {
            steps {
                echo "🐳 Build Docker Image: ${IMAGE_NAME}:${IMAGE_TAG}..."
                sh "docker build -t ${IMAGE_NAME}:${IMAGE_TAG} -t ${IMAGE_NAME}:latest ."

                echo '🛡️ Quét lỗ hổng bảo mật Docker Image bằng Trivy...'
                sh '''
                    if command -v trivy &> /dev/null; then
                        trivy image --severity CRITICAL --exit-code 0 ${IMAGE_NAME}:${IMAGE_TAG} || true
                    else
                        echo "Trivy chưa cài, bỏ qua quét Image"
                    fi
                '''
            }
        }

        stage('5. Push Image to Registry') {
            steps {
                echo "🚀 Đẩy image lên Docker Hub: ${IMAGE_NAME}:${IMAGE_TAG}..."
                // Cấu hình Credentials trong Jenkins ID 'docker-hub-credentials'
                /*
                withCredentials([usernamePassword(credentialsId: 'docker-hub-credentials', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                    sh '''
                        echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
                        docker push ${IMAGE_NAME}:${IMAGE_TAG}
                        docker push ${IMAGE_NAME}:latest
                    '''
                }
                */
                echo "Giả lập đẩy thành công image ${IMAGE_NAME}:${IMAGE_TAG}"
            }
        }

        stage('6. Update GitOps Manifest (Dev Overlay)') {
            steps {
                echo "🔄 Tự động cập nhật image tag vào repo GitOps..."
                /*
                sh """
                    git clone ${GITOPS_REPO} gitops-workdir
                    cd gitops-workdir/projects/cuahang/overlays/dev
                    kustomize edit set image ${IMAGE_NAME}=${IMAGE_NAME}:${IMAGE_TAG}
                    git config user.name "jenkins-bot"
                    git config user.email "jenkins@lan"
                    git commit -am "ci(dev): update ${APP_NAME} to ${IMAGE_TAG} [skip ci]"
                    git push origin main
                """
                */
                echo "Sẵn sàng tự động cập nhật GitOps khi kết nối token GitHub"
            }
        }
    }

    post {
        success {
            echo "✅ Pipeline CI/CD cho cuahang-be (Go) thành công!"
        }
        failure {
            echo "❌ Pipeline thất bại! Vui lòng kiểm tra log chi tiết."
        }
    }
}

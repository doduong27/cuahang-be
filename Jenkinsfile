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
            description: '🏷️ Phiên bản SemVer gốc (Dev & Stag tự gắn thêm -b<build_number>; Prod tự chốt tag sạch vx.x.x):'
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
                    def rawTag = params.IMAGE_TAG ? params.IMAGE_TAG.trim() : 'v1.0.0'
                    def cleanMatch = (rawTag =~ /^(v\d+\.\d+\.\d+)/)
                    def baseSemver = cleanMatch ? cleanMatch[0][1] : rawTag

                    if (params.TARGET_ENV == 'prod') {
                        // PRODUCTION: Luôn luôn dùng Tag sạch chuẩn quốc tế vx.x.x (ví dụ v1.0.0)
                        env.FINAL_TAG = baseSemver
                        env.IS_PROD = 'true'
                        env.SOURCE_TAG = (rawTag != baseSemver) ? rawTag : ''
                    } else if (params.TARGET_ENV == 'stag') {
                        // STAGING:
                        env.IS_PROD = 'false'
                        if (params.ACTION == 'Build & Deploy') {
                            env.FINAL_TAG = "${baseSemver}-b${BUILD_NUMBER}"
                        } else {
                            // Deploy Only: Nếu không chỉ định đuôi -b, tự tìm tag từ Dev
                            if (rawTag == baseSemver) {
                                sh """
                                    rm -rf gitops-dev-check
                                    GIT_SSH_COMMAND="ssh -o StrictHostKeyChecking=no" git clone --depth 1 ${GITOPS_REPO_SSH} gitops-dev-check
                                    DEV_TAG=\$(grep -A 2 'name: ${APP_NAME}' gitops-dev-check/projects/cuahang/overlays/dev/kustomization.yaml | grep 'newTag:' | awk '{print \$2}' | tr -d ' ' || true)
                                    rm -rf gitops-dev-check
                                    if [ -n "\$DEV_TAG" ]; then
                                        echo "\$DEV_TAG" > dev_tag.txt
                                    fi
                                """
                                if (fileExists('dev_tag.txt')) {
                                    def resolved = readFile('dev_tag.txt').trim()
                                    env.FINAL_TAG = (resolved != '') ? resolved : rawTag
                                    sh "rm -f dev_tag.txt"
                                } else {
                                    env.FINAL_TAG = rawTag
                                }
                            } else {
                                env.FINAL_TAG = rawTag
                            }
                        }
                    } else {
                        // DEV:
                        env.IS_PROD = 'false'
                        if (params.ACTION == 'Build & Deploy') {
                            env.FINAL_TAG = "${baseSemver}-b${BUILD_NUMBER}"
                        } else {
                            env.FINAL_TAG = rawTag
                        }
                    }

                    echo "=========================================================="
                    echo "🚀 KẾ HOẠCH TRIỂN KHAI BACKEND (Go Gin + PostgreSQL):"
                    echo "   - Ứng dụng: ${APP_NAME}"
                    echo "   - Môi trường đích: [ ${params.TARGET_ENV.toUpperCase()} ]"
                    echo "   - Hành động: ${params.ACTION}"
                    echo "   - Phiên bản SemVer gốc: ${baseSemver}"
                    echo "   - Image Tag triển khai: [ ${env.FINAL_TAG} ]"
                    if (env.IS_PROD == 'true') {
                        echo "   - Quy chuẩn Production: Tự động chốt Tag sạch SemVer [ ${env.FINAL_TAG} ]"
                    } else {
                        echo "   - Quy chuẩn Dev/Stag: Tự sinh Tag bất biến [ ${env.FINAL_TAG} ] (tránh cache K8s)"
                    }
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
                sh "docker build -t ${IMAGE_NAME}:${env.FINAL_TAG} -t ${IMAGE_NAME}:${params.TARGET_ENV}-latest ."

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
                        docker push ${IMAGE_NAME}:${params.TARGET_ENV}-latest
                        if [ "${env.IS_PROD}" = "true" ]; then
                            docker tag ${IMAGE_NAME}:${env.FINAL_TAG} ${IMAGE_NAME}:latest
                            docker push ${IMAGE_NAME}:latest
                        fi
                    """
                }
            }
        }

        stage('3b. Đóng dấu Tag chính thức cho Production') {
            when {
                allOf {
                    expression { return params.TARGET_ENV == 'prod' }
                    expression { return params.ACTION == 'Deploy Only (Promote)' }
                }
            }
            steps {
                script {
                    echo "📦 Đang chuẩn bị đóng dấu Tag chính thức [ ${env.FINAL_TAG} ] cho Production..."
                    withCredentials([usernamePassword(credentialsId: 'docker-hub-credentials', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                        sh """
                            echo "\$DOCKER_PASS" | docker login -u "\$DOCKER_USER" --password-stdin
                            
                            SRC_TAG="${env.SOURCE_TAG}"
                            if [ -z "\$SRC_TAG" ]; then
                                rm -rf gitops-src-check
                                GIT_SSH_COMMAND="ssh -o StrictHostKeyChecking=no" git clone --depth 1 ${GITOPS_REPO_SSH} gitops-src-check
                                SRC_TAG=\$(grep -A 2 'name: ${APP_NAME}' gitops-src-check/projects/cuahang/overlays/stag/kustomization.yaml | grep 'newTag:' | awk '{print \$2}' | tr -d ' ' || true)
                                rm -rf gitops-src-check
                            fi
                            
                            if [ -n "\$SRC_TAG" ] && [ "\$SRC_TAG" != "${env.FINAL_TAG}" ]; then
                                echo "🔄 Lấy Image đã nghiệm thu từ Staging: \${SRC_TAG} -> Gán nhãn thành Tag Prod sạch: ${env.FINAL_TAG}"
                                docker pull ${IMAGE_NAME}:\${SRC_TAG}
                                docker tag ${IMAGE_NAME}:\${SRC_TAG} ${IMAGE_NAME}:${env.FINAL_TAG}
                                docker tag ${IMAGE_NAME}:\${SRC_TAG} ${IMAGE_NAME}:latest
                                docker push ${IMAGE_NAME}:${env.FINAL_TAG}
                                docker push ${IMAGE_NAME}:latest
                                echo "✅ Đã gán nhãn và đẩy Tag Production sạch ${env.FINAL_TAG} thành công!"
                            else
                                echo "ℹ️ Image ${IMAGE_NAME}:${env.FINAL_TAG} đã sẵn sàng."
                            fi
                        """
                    }
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
                    
                    # Cập nhật image tag chuẩn xác bằng Kustomize CLI
                    kustomize edit set image ${APP_NAME}=docker.io/${IMAGE_NAME}:${env.FINAL_TAG}
                    
                    git config user.name "jenkins-bot"
                    git config user.email "jenkins@lan"
                    
                    if git status --porcelain | grep -q 'kustomization.yaml'; then
                        git add kustomization.yaml
                        git commit -m "ci(${params.TARGET_ENV}): promote ${APP_NAME} to ${env.FINAL_TAG} via Jenkins [skip ci]"
                        GIT_SSH_COMMAND="ssh -o StrictHostKeyChecking=no" git push origin main
                        echo "✅ Đã commit và push cập nhật GitOps thành công!"
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
                    def targetHost = (params.TARGET_ENV == 'dev') ? 'dev.cuahang.lan' : ((params.TARGET_ENV == 'stag') ? 'stag.cuahang.lan' : 'cuahang.lan')
                    echo "🌐 Endpoint kiểm tra API Health: http://${targetHost}/healthz"
                    echo "🌐 Endpoint danh sách sản phẩm: http://${targetHost}/api/products"
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

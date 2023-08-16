pipeline {
  agent any
  stages {
    stage('Clone dtm cluster') {
      steps {
        git(url: scm.userRemoteConfigs[0].url, branch: '$BRANCH_NAME', changelog: true, credentialsId: 'KK-github-key', poll: true)
      }
    }

    stage('Generate docker image for feature') {
      when {
        expression { BUILD_TARGET == 'true' }
        expression { BRANCH_NAME != 'master' }
      }
      steps {
        sh(returnStdout: false, script: '''
          feature_name=`echo $BRANCH_NAME | awk -F '/' '{ print $2 }'`
          mkdir -p .docker-tmp
          cp /usr/bin/consul .docker-tmp
          docker build -t uhub.service.ucloud.cn/entropypool/dtm:$feature_name .
        '''.stripIndent())
      }
    }

    stage('Build dtm image') {
      when {
        expression { BUILD_TARGET == 'true' }
        expression { BRANCH_NAME == 'master' }
      }
      steps {
        sh '''
          tag=1.17.1.1
          mkdir -p .docker-tmp
          cp /usr/bin/consul .docker-tmp
          docker build -t uhub.service.ucloud.cn/entropypool/dtm:$tag .
        '''
      }
    }

    stage('Release docker image for feature') {
      when {
        expression { RELEASE_TARGET == 'true' }
        expression { BRANCH_NAME != 'master' }
      }
      steps {
        sh(returnStdout: false, script: '''
          feature_name=`echo $BRANCH_NAME | awk -F '/' '{ print $2 }'`
          set +e
          docker images | grep dtm | grep $feature_name
          rc=$?
          set -e
          if [ 0 -eq $rc ]; then
            docker push uhub.service.ucloud.cn/entropypool/dtm:$feature_name
          fi
          images=`docker images | grep entropypool | grep dtm | grep none | awk '{ print $3 }'`
          for image in $images; do
            docker rmi $image -f
          done
        '''.stripIndent())
      }
    }

    stage('Release docker image') {
      when {
        expression { RELEASE_TARGET == 'true' }
      }
      steps {
        sh(returnStdout: true, script: '''
          tag=1.17.1.1
          set +e
          docker images | grep dtm | grep $tag
          rc=$?
          set -e
          if [ 0 -eq $rc ]; then
            docker push uhub.service.ucloud.cn/entropypool/dtm:$tag
          fi
          images=`docker images | grep entropypool | grep dtm | grep none | awk '{ print $3 }'`
          for image in $images; do
            docker rmi $image -f
          done
        '''.stripIndent())
      }
    }

    stage('Switch to current cluster') {
      steps {
        sh 'cd /etc/kubeasz; ./ezctl checkout $TARGET_ENV'
      }
    }

    stage('Deploy dtm cluster for feature') {
      when {
        expression { DEPLOY_TARGET == 'true' }
        expression { BRANCH_NAME != 'master' }
      }
      steps {
        sh(returnStdout: false, script: '''
        feature_name=`echo $BRANCH_NAME | awk -F '/' '{ print $2 }'`
        sed -ri 's#image: uhub.service.ucloud.cn/entropypool/dtm(.*)#image: uhub.service.ucloud.cn/entropypool/dtm:jenkinsfile-update#g' k8s/01-deployment.yaml
        sed -i "s/dtm.development.npool.top/dtm.$TARGET_ENV.npool.top/g" ./k8s/03-traefik-vpn-ingress.yaml
        kubectl apply -k k8s
        '''.stripIndent())
      }
    }
    stage('Deploy dtm cluster') {
      when {
        expression { DEPLOY_TARGET == 'true' }
        expression { BRANCH_NAME == 'master' }
      }
      steps {
        sh 'sed -i "s/dtm.development.npool.top/dtm.$TARGET_ENV.npool.top/g" ./k8s/03-traefik-vpn-ingress.yaml'
        sh 'kubectl apply -k k8s'
      }
    }

  }

  post('Report') {
    fixed {
      script {
        sh(script: 'bash $JENKINS_HOME/wechat-templates/send_wxmsg.sh fixed')
     }
      script {
        // env.ForEmailPlugin = env.WORKSPACE
        emailext attachmentsPattern: 'TestResults\\*.trx',
        body: '${FILE,path="$JENKINS_HOME/email-templates/success_email_tmp.html"}',
        mimeType: 'text/html',
        subject: currentBuild.currentResult + " : " + env.JOB_NAME,
        to: '$DEFAULT_RECIPIENTS'
      }
     }
    success {
      script {
        sh(script: 'bash $JENKINS_HOME/wechat-templates/send_wxmsg.sh successful')
     }
      script {
        // env.ForEmailPlugin = env.WORKSPACE
        emailext attachmentsPattern: 'TestResults\\*.trx',
        body: '${FILE,path="$JENKINS_HOME/email-templates/success_email_tmp.html"}',
        mimeType: 'text/html',
        subject: currentBuild.currentResult + " : " + env.JOB_NAME,
        to: '$DEFAULT_RECIPIENTS'
      }
     }
    failure {
      script {
        sh(script: 'bash $JENKINS_HOME/wechat-templates/send_wxmsg.sh failure')
     }
      script {
        // env.ForEmailPlugin = env.WORKSPACE
        emailext attachmentsPattern: 'TestResults\\*.trx',
        body: '${FILE,path="$JENKINS_HOME/email-templates/fail_email_tmp.html"}',
        mimeType: 'text/html',
        subject: currentBuild.currentResult + " : " + env.JOB_NAME,
        to: '$DEFAULT_RECIPIENTS'
      }
     }
    aborted {
      script {
        sh(script: 'bash $JENKINS_HOME/wechat-templates/send_wxmsg.sh aborted')
     }
      script {
        // env.ForEmailPlugin = env.WORKSPACE
        emailext attachmentsPattern: 'TestResults\\*.trx',
        body: '${FILE,path="$JENKINS_HOME/email-templates/fail_email_tmp.html"}',
        mimeType: 'text/html',
        subject: currentBuild.currentResult + " : " + env.JOB_NAME,
        to: '$DEFAULT_RECIPIENTS'
      }
     }
  }
}

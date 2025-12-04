# Tiltfile

# Define the docker images for each service
docker_build('auth-service', './', dockerfile='services/auth/Dockerfile')
docker_build('order-service', './', dockerfile='services/order/Dockerfile')
docker_build('shipping-service', './', dockerfile='services/shipping/Dockerfile')
docker_build('notification-service', './', dockerfile='services/notification/Dockerfile')
docker_build('analytics-service', './', dockerfile='services/analytics/Dockerfile')
docker_build('bff-service', './', dockerfile='services/bff/Dockerfile')
docker_build('frontend', './', dockerfile='frontend/Dockerfile')

# Deploy using Helm
k8s_yaml(helm(
    './deploy/helm/my-store',
    name='my-store',
    values=['./deploy/helm/my-store/values.yaml'],
    set=[
        # Ensure we use the images we built
        'services.auth.image.repository=auth-service',
        'services.order.image.repository=order-service',
        'services.shipping.image.repository=shipping-service',
        'services.notification.image.repository=notification-service',
        'services.analytics.image.repository=analytics-service',
        'services.bff.image.repository=bff-service',
        'services.frontend.image.repository=frontend',
        # Set pull policy to ensure we use local images
        'services.auth.image.pullPolicy=IfNotPresent',
        'services.order.image.pullPolicy=IfNotPresent',
        'services.shipping.image.pullPolicy=IfNotPresent',
        'services.notification.image.pullPolicy=IfNotPresent',
        'services.analytics.image.pullPolicy=IfNotPresent',
        'services.bff.image.pullPolicy=IfNotPresent',
        'services.frontend.image.pullPolicy=IfNotPresent',
    ]
))

# Configure resources (Port forwarding and labels)
# Note: Resource names match the deployment names in the Helm chart

k8s_resource('frontend', 
    port_forwards='3000:80',
    labels=['frontend']
)

k8s_resource('bff', 
    port_forwards='8080:8080',
    labels=['backend']
)

# Group backend services for cleaner UI
backend_services = ['auth', 'order', 'shipping', 'notification', 'analytics']
for service in backend_services:
    # Using .format() instead of f-strings for compatibility
    k8s_resource(service, labels=['backend'])

# Infrastructure resources
k8s_resource('postgres', labels=['infra'])
k8s_resource('pgadmin', 
    port_forwards='5050:80',
    labels=['infra']
)
k8s_resource('kafka', labels=['infra'])
k8s_resource('zookeeper', labels=['infra'])

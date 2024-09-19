/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

export default {
  'component.upstream.fields.tls.client_key': 'Clé du client',
  'component.upstream.fields.tls.client_key.required': 'Veuillez entrer la clé du client',
  'component.upstream.fields.tls.client_cert': 'Certificat du client',
  'component.upstream.fields.tls.client_cert.required': 'Veuillez entrer le certificat du client',
  'component.upstream.fields.upstream_type': "Type d'amont",
  'component.upstream.fields.upstream_type.node': 'Noeud',
  'component.upstream.fields.upstream_type.service_discovery': 'Découverte de service',
  'component.upstream.fields.discovery_type': 'Type de découverte',
  'component.upstream.fields.discovery_type.tooltip': 'Type de découverte',
  'component.upstream.fields.discovery_type.placeholder':
    'Veuillez sélectionner le type de découverte',
  'component.upstream.fields.discovery_type.type.dns': 'DNS',
  'component.upstream.fields.discovery_type.type.consul': 'Consul',
  'component.upstream.fields.discovery_type.type.consul_kv': 'Consul KV',
  'component.upstream.fields.discovery_type.type.nacos': 'Nacos',
  'component.upstream.fields.discovery_type.type.eureka': 'Eureka',
  'component.upstream.fields.discovery_type.type.kubernetes': 'Kubernetes',
  'component.upstream.fields.discovery_args.group_name': 'Nom du groupe',
  'component.upstream.fields.discovery_args.group_name.tooltip': 'Nom du groupe',
  'component.upstream.fields.discovery_args.group_name.placeholder':
    'Veuillez entrer le nom du groupe',
  'component.upstream.fields.discovery_args.namespace_id': "ID de l'espace de noms",
  'component.upstream.fields.discovery_args.namespace_id.tooltip': "ID de l'espace de noms",
  'component.upstream.fields.discovery_args.namespace_id.placeholder':
    "Veuillez entrer l'ID de l'espace de noms",
  'component.upstream.fields.service_name': 'Nom du service',
  'component.upstream.fields.service_name.tooltip': 'Nom du service',
  'component.upstream.fields.service_name.placeholder': 'Veuillez entrer le nom du service',
  'component.upstream.fields.scheme.tooltip.stream':
    'Ce type est utilisé uniquement pour Stream Route, qui est un proxy de la couche 4. Référence: https://apisix.apache.org/docs/apisix/stream-proxy/',
  'component.upstream.fields.scheme.tooltip.pubsub':
    'Ce type est utilisé uniquement dans la publication-abonnement. Référence: https://apisix.apache.org/docs/apisix/pubsub/',
  'component.upstream.fields.tls': 'TLS',
  'component.upstream.fields.tls.tooltip': 'Certificat TLS',
  'component.upstream.fields.hash_on': 'Hash sur',
  'component.upstream.fields.hash_on.tooltip': 'Ce qui doit être utilisé comme entrée de hachage',
  'component.upstream.fields.key': 'Clé',
  'component.upstream.fields.key.tooltip': "Clé en tant qu'entrée de hachage",
  'component.upstream.fields.retries': 'Réessais',
  'component.upstream.fields.retries.tooltip':
    'Le mécanisme de réessai envoie la demande au nœud amont suivant. Une valeur de 0 désactive le mécanisme de réessai et laisse la table vide pour utiliser le nombre de nœuds backend disponibles.',
  'component.upstream.fields.retry_timeout': 'Délai de réessai',
  'component.upstream.fields.retry_timeout.tooltip':
    'Configurer un nombre pour limiter la durée en secondes pendant laquelle les réessais peuvent se poursuivre, et ne pas continuer les réessais si la demande précédente et les réessais ont pris trop de temps. 0 signifie désactiver le mécanisme de délai de réessai.',
  'component.upstream.fields.keepalive_pool': 'Pool Keepalive',
  'component.upstream.fields.keepalive_pool.tooltip':
    "Définir un pool Keepalive indépendant pour l'amont",
  'component.upstream.fields.keepalive_pool.size': 'Taille',
  'component.upstream.fields.keepalive_pool.size.placeholder': 'Veuillez entrer la taille',
  'component.upstream.fields.keepalive_pool.idle_timeout': "Délai d'inactivité",
  'component.upstream.fields.keepalive_pool.idle_timeout.placeholder':
    "Veuillez entrer le délai d'inactivité",
  'component.upstream.fields.keepalive_pool.requests': 'Demandes',
  'component.upstream.fields.keepalive_pool.requests.placeholder': 'Veuillez entrer les demandes',
  'component.upstream.fields.checks.active.type': 'Type',
  'component.upstream.fields.checks.active.type.tooltip':
    "Indiquez s'il faut effectuer des vérifications de santé actives à l'aide d'HTTP ou HTTPS, ou simplement tenter une connexion TCP.",
  'component.upstream.fields.checks.active.concurrency': 'Concurrence',
  'component.upstream.fields.checks.active.concurrency.tooltip':
    'Nombre de cibles à vérifier simultanément dans les vérifications de santé actives.',
  'component.upstream.fields.checks.active.host': 'Hôte',
  'component.upstream.fields.checks.active.host.required': "Veuillez entrer le nom d'hôte",
  'component.upstream.fields.checks.active.host.tooltip':
    "Le nom d'hôte de la requête HTTP utilisé pour effectuer la vérification de santé active",
  'component.upstream.fields.checks.active.host.scope':
    'Seules les lettres, les chiffres et . sont pris en charge',
  'component.upstream.fields.checks.active.port': 'Port',
  'component.upstream.fields.checks.active.http_path': 'Chemin HTTP',
  'component.upstream.fields.checks.active.http_path.tooltip':
    "Le chemin qui doit être utilisé lors de l'émission de la requête GET HTTP à la cible. La valeur par défaut est /.",
  'component.upstream.fields.checks.active.http_path.placeholder':
    'Veuillez entrer le chemin de la requête HTTP',
  'component.upstream.fields.checks.active.https_verify_certificate':
    'Vérifier le certificat HTTPs',
  'component.upstream.fields.checks.active.https_verify_certificate.tooltip':
    "Indique s'il faut vérifier la validité du certificat SSL de l'hôte distant lors de la réalisation de vérifications de santé actives à l'aide de HTTPS.",
  'component.upstream.fields.checks.active.req_headers': 'En-têtes de la demande',
  'component.upstream.fields.checks.active.req_headers.tooltip':
    'En-têtes de demande supplémentaires, exemple: User-Agent: curl/7.29.0',
  'component.upstream.fields.checks.active.healthy.interval': 'Intervalle',
  'component.upstream.fields.checks.active.healthy.interval.tooltip':
    'Intervalle entre les vérifications des cibles saines (en secondes)',
  'component.upstream.fields.checks.active.healthy.successes': 'Réussites',
  'component.upstream.fields.checks.active.healthy.successes.tooltip':
    'Nombre de succès pour considérer une cible comme saine',
  'component.upstream.fields.checks.active.healthy.successes.required':
    'Veuillez entrer le nombre de réussites',
  'component.upstream.fields.checks.active.healthy.http_statuses': 'Statuts HTTP',
  'component.upstream.fields.checks.active.healthy.http_statuses.tooltip':
    "Un tableau de statuts HTTP à considérer comme un succès, indiquant la santé, lorsqu'il est renvoyé par une sonde dans les vérifications de santé actives.",
  'component.upstream.fields.checks.active.unhealthy.timeouts': 'Délais dépassés',
  'component.upstream.fields.checks.active.unhealthy.timeouts.tooltip':
    'Nombre de délais dépassés dans les sondes actives pour considérer une cible comme non saine.',
  'component.upstream.fields.checks.active.unhealthy.http_failures': 'Échecs HTTP',
  'component.upstream.fields.checks.active.unhealthy.http_failures.tooltip':
    "Nombre d'échecs HTTP pour considérer une cible comme non saine",
  'component.upstream.fields.checks.active.unhealthy.http_failures.required':
    'Veuillez entrer les échecs HTTP',
  'component.upstream.fields.checks.active.unhealthy.tcp_failures': 'Échecs TCP',
  'component.upstream.fields.checks.active.unhealthy.tcp_failures.tooltip':
    "Nombre d'échecs TCP pour considérer une cible comme non saine",
  'component.upstream.fields.checks.active.unhealthy.tcp_failures.required':
    'Veuillez entrer les échecs TCP',
  'component.upstream.fields.checks.active.unhealthy.interval': 'Intervalle',
  'component.upstream.fields.checks.active.unhealthy.interval.tooltip':
    'Intervalle entre les vérifications de santé actives pour les cibles non saines (en secondes). Une valeur de zéro indique que les sondes actives pour les cibles saines ne doivent pas être effectuées.',
  'component.upstream.fields.checks.active.unhealthy.required':
    "Veuillez entrer l'intervalle non sain",
  'component.upstream.fields.checks.passive.healthy.successes': 'Réussites',
  'component.upstream.fields.checks.passive.healthy.successes.tooltip':
    'Nombre de succès pour considérer une cible comme saine',
  'component.upstream.fields.checks.passive.healthy.successes.required':
    'Veuillez entrer le nombre de réussites',
  'component.upstream.fields.checks.passive.unhealthy.timeouts': 'Délais dépassés',
  'component.upstream.fields.checks.passive.unhealthy.timeouts.tooltip':
    "Nombre de délais dépassés dans le trafic proxy pour considérer une cible comme non saine, tel qu'observé par les vérifications de santé passives.",
  'component.upstream.other.none': 'Aucun (Disponible uniquement lors de la liaison du service)',
  'component.upstream.other.pass_host-with-multiple-nodes.title':
    'Veuillez vérifier la configuration du nœud cible',
  'component.upstream.other.pass_host-with-multiple-nodes':
    "Lors de l'utilisation d'un nom d'hôte ou d'une adresse IP dans la liste des nœuds cibles, assurez-vous qu'il n'y a qu'un seul nœud cible",
  'component.upstream.other.health-check.passive-only':
    'Lorsque la vérification de santé passive est activée, la vérification de santé active doit également être activée.',
  'component.upstream.other.health-check.invalid':
    'Veuillez vérifier la configuration de la vérification de santé',
};

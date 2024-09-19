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
  'page.upstream.step.select.upstream': 'Amont',
  'page.upstream.step.select.upstream.select.option': 'Personnalisé',
  'page.upstream.step.select.upstream.select.option.serviceSelected':
    'Personnalisé (La configuration actuelle remplacera le service lié)',
  'page.upstream.step.select.upstream.select.none': 'Aucun',
  'page.upstream.step.backend.server.domain.or.ip': 'Hôte/IP du serveur backend',
  'page.upstream.form.item-label.node.domain.or.ip': 'Cibles',
  'page.upstream.step.input.domain.name.or.ip': "Veuillez entrer le domaine ou l'IP",
  'page.upstream.step.valid.domain.name.or.ip': 'Veuillez entrer un domaine ou une IP valide',
  'page.upstream.step.domain.name.or.ip': "Nom d'hôte ou IP",
  'page.upstream.step.host': 'Hôte',
  'page.upstream.step.input.port': 'Veuillez entrer le numéro de port',
  'page.upstream.step.port': 'Port',
  'page.upstream.step.input.weight': 'Veuillez entrer le poids',
  'page.upstream.step.weight': 'Poids',
  'page.upstream.step.create': 'Créer',
  'page.upstream.step.name': 'Nom',
  'page.upstream.step.name.should.unique': 'Le nom doit être unique',
  'page.upstream.step.input.upstream.name': "Veuillez entrer le nom de l'amont",
  'page.upstream.step.description': 'Description',
  'page.upstream.step.input.description': "Veuillez entrer la description de l'amont",
  'page.upstream.step.type': 'Algorithme',
  'page.upstream.step.pass-host': "Nom d'hôte",
  'page.upstream.step.pass-host.pass': 'Conserver le même hôte de la demande client',
  'page.upstream.step.pass-host.node': "Utiliser le domaine ou l'IP de la liste des nœuds",
  'page.upstream.step.pass-host.rewrite': "Hôte personnalisé (Sera obsolète à l'avenir)",
  'page.upstream.step.pass-host.upstream_host': 'Hôte personnalisé',
  'page.upstream.step.connect.timeout': 'Délai de connexion',
  'page.upstream.step.connect.timeout.desc':
    "Délai d'établissement d'une connexion de la demande au serveur amont",
  'page.upstream.step.input.connect.timeout': 'Veuillez entrer le délai de connexion',
  'page.upstream.step.send.timeout': "Délai d'envoi",
  'page.upstream.step.send.timeout.desc': "Délai d'envoi des données aux serveurs amont",
  'page.upstream.step.input.send.timeout': "Veuillez entrer le délai d'envoi",
  'page.upstream.step.read.timeout': 'Délai de lecture',
  'page.upstream.step.read.timeout.desc': 'Délai de réception des données des serveurs amont',
  'page.upstream.step.input.read.timeout': 'Veuillez entrer le délai de lecture',
  'page.upstream.step.healthyCheck.healthy.check': 'Vérification de santé',
  'page.upstream.step.healthyCheck.healthy': 'Sain',
  'page.upstream.step.healthyCheck.unhealthy': 'Malsain',
  'page.upstream.step.healthyCheck.healthy.status': 'État sain',
  'page.upstream.step.healthyCheck.unhealthyStatus': 'État malsain',
  'page.upstream.step.healthyCheck.active': 'Actif',
  'page.upstream.step.healthyCheck.active.timeout': 'Délai',
  'page.upstream.step.input.healthyCheck.active.timeout': 'Veuillez entrer le délai',
  'page.upstream.step.input.healthyCheck.activeInterval': "Veuillez entrer l'intervalle",
  'page.upstream.step.input.healthyCheck.active.req_headers':
    'Veuillez entrer les en-têtes de demande',
  'page.upstream.step.healthyCheck.passive': 'Passif',
  'page.upstream.step.healthyCheck.passive.http_statuses': 'Statuts HTTP',
  'page.upstream.step.input.healthyCheck.passive.http_statuses': 'Veuillez entrer le statut HTTP',
  'page.upstream.step.healthyCheck.passive.tcp_failures': 'Échecs TCP',
  'page.upstream.step.input.healthyCheck.passive.tcp_failures': 'Veuillez entrer les échecs TCP',
  'page.upstream.step.keepalive_pool': 'Pool de maintien de la connexion',
  'page.upstream.notificationMessage.enableHealthCheckFirst':
    "Veuillez d'abord activer la vérification de santé.",
  'page.upstream.upstream_host.required': "Veuillez entrer l'hôte personnalisé",
  'page.upstream.create': "Créer l'amont",
  'page.upstream.configure': "Configurer l'amont",
  'page.upstream.create.upstream.successfully': "Création de l'amont réussie",
  'page.upstream.edit.upstream.successfully': "Configuration de l'amont réussie",
  'page.upstream.create.basic.info': 'Informations de base',
  'page.upstream.create.preview': 'Aperçu',
  'page.upstream.list.id': 'ID',
  'page.upstream.list.name': 'Nom',
  'page.upstream.list.type': 'Type',
  'page.upstream.list.description': 'Description',
  'page.upstream.list.edit.time': 'Heure de mise à jour',
  'page.upstream.list.operation': 'Opération',
  'page.upstream.list.edit': 'Configurer',
  'page.upstream.list.confirm.delete': 'Êtes-vous sûr de vouloir supprimer ?',
  'page.upstream.list.confirm': 'Confirmer',
  'page.upstream.list.cancel': 'Annuler',
  'page.upstream.list.delete.successfully': "Suppression de l'amont réussie",
  'page.upstream.list.delete': 'Supprimer',
  'page.upstream.list': 'Liste des amonts',
  'page.upstream.list.input': 'Veuillez entrer',
  'page.upstream.list.create': 'Créer',
  'page.upstream.type.roundrobin': 'Round Robin',
  'page.upstream.type.chash': 'CHash',
  'page.upstream.type.ewma': 'EWMA',
  'page.upstream.type.least_conn': 'Moins de conn',
  'page.upstream.list.content':
    "La liste des amonts contient les services amonts créés (c'est-à-dire, les services backend) et permet l'équilibrage de charge et la vérification de l'état de plusieurs nœuds cibles des services amonts.",
  'page.upstream.checks.active.timeout.description':
    'Délai de socket pour les vérifications actives (en secondes)',
  'page.upstream.checks.active.unhealthy.interval.description':
    'Intervalle entre les vérifications des cibles malsaines (en secondes)',
  'page.upstream.checks.passive.healthy.http_statuses.description':
    'Quels statuts HTTP considérer comme un succès',
  'page.upstream.checks.passive.unhealthy.http_statuses.description':
    'Quels statuts HTTP considérer comme un échec',
  'page.upstream.checks.passive.unhealthy.http_failures.description':
    "Nombre d'échecs HTTP à considérer comme une cible malsaine",
  'page.upstream.checks.passive.unhealthy.tcp_failures.description':
    "Nombre d'échecs TCP à considérer comme une cible malsaine",
  'page.upstream.scheme': 'Schéma',
  'page.upstream.other.configuration.invalid': "Veuillez vérifier la configuration de l'amont",
};

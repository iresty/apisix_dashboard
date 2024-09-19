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
  'component.plugin.tip1':
    'NOTE : Après avoir personnalisé le plugin, vous devez mettre à jour schema.json.',
  'component.plugin.tip2': 'Comment mettre à jour ?',
  'component.select.pluginTemplate': 'Sélectionnez un modèle de plugin',
  'component.step.select.pluginTemplate.select.option': 'Personnalisé',
  'component.plugin.pluginTemplate.tip1':
    "1. Lorsqu'une route a déjà le champ plugins configuré, les plugins dans le modèle de plugin seront fusionnés avec celui-ci.",
  'component.plugin.pluginTemplate.tip2':
    '2. Le même plugin dans le modèle de plugin remplacera celui dans les plugins.',
  'component.plugin.enable': 'Activer',
  'component.plugin.disable': 'Modifier',
  'component.plugin.authentication': 'Authentification',
  'component.plugin.security': 'Sécurité',
  'component.plugin.traffic': 'Contrôle du trafic',
  'component.plugin.serverless': 'Serverless',
  'component.plugin.observability': 'Observabilité',
  'component.plugin.other': 'Autre',
  'component.plugin.all': 'Tout',
  // cors
  'component.pluginForm.cors.allow_origins.tooltip':
    "Quelles origines sont autorisées à activer CORS, au format : scheme://host:port, par exemple : https://somehost.com:8081. Utilisez , pour séparer plusieurs origines. Lorsque allow_credential est false, vous pouvez utiliser * pour indiquer d'autoriser n'importe quelle origine. Vous pouvez également autoriser toutes les origines de force en utilisant ** même si allow_credential est déjà activé, mais cela présentera des risques de sécurité.",
  'component.pluginForm.cors.allow_origins.extra': 'Par exemple : https://somehost.com:8081',
  'component.pluginForm.cors.allow_methods.tooltip':
    "Quelles méthodes sont autorisées à activer CORS, telles que : GET, POST, etc. Utilisez , pour séparer plusieurs méthodes. Lorsque allow_credential est false, vous pouvez utiliser * pour indiquer d'autoriser toutes les méthodes. Vous pouvez également autoriser n'importe quelle méthode de force en utilisant ** même si allow_credential est déjà activé, mais cela présentera des risques de sécurité.",
  'component.pluginForm.cors.allow_headers.tooltip':
    "Quels en-têtes sont autorisés à définir dans la demande lors de l'accès à une ressource en cross-origin. Utilisez , pour séparer plusieurs valeurs. Lorsque allow_credential est false, vous pouvez utiliser * pour indiquer d'autoriser tous les en-têtes de demande. Vous pouvez également autoriser n'importe quel en-tête de force en utilisant ** même si allow_credential est déjà activé, mais cela présentera des risques de sécurité.",
  'component.pluginForm.cors.expose_headers.tooltip':
    "Quels en-têtes sont autorisés à définir dans la réponse lors de l'accès à une ressource en cross-origin. Utilisez , pour séparer plusieurs valeurs. Lorsque allow_credential est false, vous pouvez utiliser * pour indiquer d'autoriser n'importe quel en-tête. Vous pouvez également autoriser n'importe quel en-tête de force en utilisant ** même si allow_credential est déjà activé, mais cela présentera des risques de sécurité.",
  'component.pluginForm.cors.max_age.tooltip':
    'Nombre maximal de secondes pendant lesquelles les résultats peuvent être mis en cache. Pendant cette période, le navigateur réutilisera le dernier résultat de vérification. -1 signifie pas de mise en cache. Veuillez noter que la valeur maximale dépend du navigateur, veuillez vous référer à MDN pour plus de détails.',
  'component.pluginForm.cors.allow_credential.tooltip':
    "Si vous définissez cette option sur true, vous ne pouvez pas utiliser '*' pour les autres options.",
  'component.pluginForm.cors.allow_origins_by_metadata.tooltip':
    'Correspondance avec quelle origine est autorisée à activer CORS en référençant allow_origins défini dans les métadonnées du plugin.',
  'component.pluginForm.cors.allow_origins_by_regex.tooltip':
    'Utilisez des expressions regex pour correspondre à quelle origine est autorisée à activer CORS. Chaque champ de saisie ne peut être configuré qu\'avec une seule expression régulière autonome, par exemple ".*.test.com" qui peut correspondre à n\'importe quel sous-domaine de test.com.',
  // referer-restriction
  'component.pluginForm.referer-restriction.whitelist.tooltip':
    "Liste des noms d'hôte à autoriser. Le nom d'hôte peut commencer par * comme joker.",
  'component.pluginForm.referer-restriction.blacklist.tooltip':
    "Liste des noms d'hôte à mettre sur liste noire. Le nom d'hôte peut commencer par * comme joker.",
  'component.pluginForm.referer-restriction.listEmpty.tooltip': 'Liste vide',
  'component.pluginForm.referer-restriction.bypass_missing.tooltip':
    "Indique s'il faut contourner la vérification lorsque l'en-tête Referer est manquant ou mal formé.",
  'component.pluginForm.referer-restriction.message.tooltip':
    "Message renvoyé en cas d'accès non autorisé.",
  // api-breaker
  'component.pluginForm.api-breaker.break_response_code.tooltip':
    "Code d'erreur renvoyé en cas de problème de santé.",
  'component.pluginForm.api-breaker.break_response_body.tooltip':
    "Corps du message de réponse lorsque l'amont est en mauvais état de santé.",
  'component.pluginForm.api-breaker.break_response_headers.tooltip':
    'En-têtes renvoyés en cas de problème de santé.',
  'component.pluginForm.api-breaker.max_breaker_sec.tooltip':
    'Temps maximal du disjoncteur (en secondes).',
  'component.pluginForm.api-breaker.unhealthy.http_statuses.tooltip':
    "Codes d'état en cas de problème de santé.",
  'component.pluginForm.api-breaker.unhealthy.failures.tooltip':
    "Nombre de demandes d'erreur consécutives qui ont déclenché un état de santé défavorable.",
  'component.pluginForm.api-breaker.healthy.http_statuses.tooltip': "Codes d'état en cas de santé.",
  'component.pluginForm.api-breaker.healthy.successes.tooltip':
    'Nombre de demandes normales consécutives qui déclenchent un état de santé.',
  // proxy-mirror
  'component.pluginForm.proxy-mirror.host.tooltip':
    "Spécifiez une adresse de service miroir, par exemple http://127.0.0.1:9797 (l'adresse doit contenir un schéma : http ou https, pas de partie URI)",
  'component.pluginForm.proxy-mirror.host.extra': 'Par exemple : http://127.0.0.1:9797',
  'component.pluginForm.proxy-mirror.host.ruletip':
    "l'adresse doit contenir un schéma : http ou https, pas de partie URI",
  'component.pluginForm.proxy-mirror.path.tooltip':
    'Spécifiez la partie chemin de la demande miroir. Sans cela, le chemin actuel sera utilisé.',
  'component.pluginForm.proxy-mirror.path.ruletip':
    'Veuillez saisir le chemin correct, par exemple /path',
  'component.pluginForm.proxy-mirror.sample_ratio.tooltip':
    "le taux d'échantillonnage pour lequel les demandes seront miroitées.",
  // limit-conn
  'component.pluginForm.limit-conn.conn.tooltip':
    'le nombre maximum de demandes simultanées autorisées. Les demandes dépassant ce ratio (et en dessous de conn + burst) seront retardées (la latence est configurée par défaut par default_conn_delay) pour se conformer à ce seuil.',
  'component.pluginForm.limit-conn.burst.tooltip':
    'le nombre de demandes simultanées excessives (ou de connexions) autorisées à être retardées.',
  'component.pluginForm.limit-conn.default_conn_delay.tooltip':
    'la latence des secondes de la demande lorsque les demandes simultanées dépassent conn mais sont en dessous de (conn + burst).',
  'component.pluginForm.limit-conn.key_type.tooltip':
    'Le type de clé, supporte : "var" (variable unique) et "var_combination" (combinaison de variable)',
  'component.pluginForm.limit-conn.key.tooltip':
    "pour limiter le niveau de concurrence. Par exemple, on peut utiliser le nom d'hôte (ou la zone du serveur) comme clé afin de limiter la concurrence par nom d'hôte. Sinon, on peut également utiliser l'adresse client comme clé afin d'éviter qu'un seul client n'inonde notre service de trop de connexions ou de demandes parallèles.",
  'component.pluginForm.limit-conn.rejected_code.tooltip':
    'retourné lorsque la demande dépasse conn + burst sera rejetée.',
  'component.pluginForm.limit-conn.rejected_msg.tooltip':
    'le corps de la réponse retourné lorsque la demande dépasse conn + burst sera rejetée.',
  'component.pluginForm.limit-conn.only_use_default_delay.tooltip':
    'activez le mode strict des secondes de latence. Si vous définissez cette option sur true, il fonctionnera strictement selon les secondes de latence que vous avez définies sans logique de calcul supplémentaire.',
  'component.pluginForm.limit-conn.allow_degradation.tooltip':
    "Indiquez s'il faut activer la dégradation du plugin lorsque la fonction limit-conn est temporairement indisponible. Autorisez les demandes à continuer lorsque la valeur est définie sur true, par défaut false.",
  // limit-req
  'component.pluginForm.limit-req.rate.tooltip':
    'Le taux de demande spécifié (nombre par seconde) seuil. Les demandes dépassant ce taux (et en dessous de burst) seront retardées pour se conformer au taux.',
  'component.pluginForm.limit-req.burst.tooltip':
    'Le nombre de demandes excessives par seconde autorisées à être retardées. Les demandes dépassant cette limite dure seront rejetées immédiatement.',
  'component.pluginForm.limit-req.key_type.tooltip':
    'Le type de clé, supporte : "var" (variable unique) et "var_combination" (combinaison de variable)',
  'component.pluginForm.limit-req.key.tooltip':
    "La clé spécifiée par l'utilisateur pour limiter le taux.",
  'component.pluginForm.limit-req.rejected_code.tooltip':
    "Le code d'état HTTP renvoyé lorsque la demande dépasse le seuil est rejeté.",
  'component.pluginForm.limit-req.rejected_msg.tooltip':
    'Le corps de la réponse renvoyé lorsque la demande dépasse le seuil est rejeté.',
  'component.pluginForm.limit-req.nodelay.tooltip':
    'Si le drapeau nodelay est true, les demandes en rafale ne seront pas retardées.',
  'component.plugin.form': 'Formulaire',
  'component.plugin.format-codes.disable': 'Formater les données JSON ou YAML',
  'component.plugin.editor': 'Éditeur de plugin',
  'component.plugin.noConfigurationRequired': "N'a pas besoin de configuration",
  // limit-count
  'component.pluginForm.limit-count.count.tooltip': 'Le nombre spécifié de seuil de demandes.',
  'component.pluginForm.limit-count.time_window.tooltip':
    'La fenêtre de temps en secondes avant que le nombre de demandes ne soit réinitialisé.',
  'component.pluginForm.limit-count.key_type.tooltip':
    'Le type de clé, supporte : "var" (variable unique) et "var_combination" (combinaison de variable)',
  'component.pluginForm.limit-count.key.tooltip':
    "La clé spécifiée par l'utilisateur pour limiter le nombre.",
  'component.pluginForm.limit-count.rejected_code.tooltip':
    "Le code d'état HTTP renvoyé lorsque la demande dépasse le seuil est rejeté, par défaut 503.",
  'component.pluginForm.limit-count.rejected_msg.tooltip':
    'Le corps de la réponse renvoyé lorsque la demande dépasse le seuil est rejeté.',
  'component.pluginForm.limit-count.policy.tooltip':
    "Les politiques de limitation de débit à utiliser pour récupérer et incrémenter les limites. Les valeurs disponibles sont local(les compteurs seront stockés localement en mémoire sur le nœud) et redis(les compteurs sont stockés sur un serveur Redis et seront partagés entre les nœuds, l'utiliser généralement pour faire la limite de vitesse globale) et redis-cluster(la même fonction que redis, utilisez uniquement le modèle de cluster Redis).",
  'component.pluginForm.limit-count.allow_degradation.tooltip':
    "Indiquez s'il faut activer la dégradation du plugin lorsque la fonction limit-count est temporairement indisponible (par exemple, délai d'attente de redis). Autorisez les demandes à continuer lorsque la valeur est définie sur true",
  'component.pluginForm.limit-count.show_limit_quota_header.tooltip':
    "Indiquez s'il faut afficher X-RateLimit-Limit et X-RateLimit-Remaining (ce qui signifie le nombre total de demandes et le nombre restant de demandes qui peuvent être envoyées) dans l'en-tête de réponse",
  'component.pluginForm.limit-count.group.tooltip':
    'La route configurée avec le même groupe partagera le même compteur',
  'component.pluginForm.limit-count.redis_host.tooltip':
    "Lors de l'utilisation de la politique redis, cette propriété spécifie l'adresse du serveur Redis.",
  'component.pluginForm.limit-count.redis_port.tooltip':
    "Lors de l'utilisation de la politique redis, cette propriété spécifie le port du serveur Redis.",
  'component.pluginForm.limit-count.redis_password.tooltip':
    "Lors de l'utilisation de la politique redis, cette propriété spécifie le mot de passe du serveur Redis.",
  'component.pluginForm.limit-count.redis_database.tooltip':
    "Lors de l'utilisation de la politique redis, cette propriété spécifie la base de données que vous avez sélectionnée du serveur Redis, et uniquement pour le mode non cluster Redis (mode instance unique ou service cloud public Redis qui fournit une seule entrée).",
  'component.pluginForm.limit-count.redis_timeout.tooltip':
    "Lors de l'utilisation de la politique redis, cette propriété spécifie le délai en millisecondes de toute commande soumise au serveur Redis.",
  'component.pluginForm.limit-count.redis_cluster_nodes.tooltip':
    "Lors de l'utilisation de la politique redis-cluster, cette propriété est une liste d'adresses des nœuds de service Redis cluster (au moins deux nœuds).",
  'component.pluginForm.limit-count.redis_cluster_name.tooltip':
    "Lors de l'utilisation de la politique redis-cluster, cette propriété est le nom des nœuds de service Redis cluster.",
  'component.pluginForm.limit-count.atLeast2Characters.rule':
    'Veuillez entrer au moins 2 caractères',
};

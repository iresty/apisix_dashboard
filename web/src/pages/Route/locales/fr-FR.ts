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
  'page.route.button.returnList': 'Aller à la liste',
  'page.route.button.send': 'Envoyer',
  'page.route.onlineDebug': 'Débogage en ligne',
  'page.route.pluginTemplateConfig': 'Configuration du modèle de plugin',
  'page.route.parameterPosition': 'Position du paramètre',
  'page.route.httpRequestHeader': 'En-tête de requête HTTP',
  'page.route.requestParameter': 'Paramètre de requête',
  'page.route.postRequestParameter': 'Paramètre de requête POST',
  'page.route.builtinParameter': 'Paramètre intégré',
  'page.route.parameterName': 'Nom du paramètre',
  'page.route.operationalCharacter': 'Caractère opérationnel',
  'page.route.equal': 'Égal(==)',
  'page.route.unequal': 'Différent(~=)',
  'page.route.greaterThan': 'Supérieur à(>)',
  'page.route.lessThan': 'Inférieur à(<)',
  'page.route.regexMatch': 'Correspondance regex(~~)',
  'page.route.caseInsensitiveRegexMatch': 'Correspondance regex insensible à la casse(~*)',
  'page.route.in': 'DANS',
  'page.route.has': 'A',
  'page.route.reverse': 'Inverser le résultat(!)',
  'page.route.rule': 'Règle',
  'page.route.httpHeaderName': "Nom de l'en-tête de requête HTTP",
  'page.route.service': 'Service',
  'page.route.instructions': 'Instructions',
  'page.route.import': 'Importer',
  'page.route.createRoute': 'Créer une route',
  'page.route.editRoute': 'Configurer la route',
  'page.route.batchDeletion': 'Suppression en lot des routes',
  'page.route.unSelect': 'Désélectionner',
  'page.route.item': 'éléments',
  'page.route.chosen': 'choisis',
  'page.route.input.placeholder.parameterNameHttpHeader':
    "Nom de l'en-tête de requête, par exemple : HOST",
  'page.route.input.placeholder.parameterNameRequestParameter':
    'Nom du paramètre, par exemple : id',
  'page.route.input.placeholder.requestUrl': "Veuillez saisir l'URL de requête valide",
  'page.route.input.placeholder.paramKey': 'Clé du paramètre',
  'page.route.input.placeholder.paramValue': 'Valeur du paramètre',
  'page.route.input.placeholder.paramType': 'Type de paramètre',
  'page.route.form.itemRulesRequiredMessage.parameterName':
    'Seules les lettres et les chiffres sont autorisés, et le nom doit commencer par une lettre',
  'page.route.value': 'Valeur du paramètre',
  'page.route.panelSection.title.advancedMatchRule':
    'Conditions avancées de correspondance des routes',
  'page.route.panelSection.title.nameDescription': 'Nom et description',
  'page.route.form.itemRulesPatternMessage.apiNameRule':
    'La longueur maximale doit être de 100 seulement',
  'page.route.panelSection.title.requestConfigBasicDefine': 'Définition de base de la requête',
  'page.route.form.itemLabel.httpMethod': 'Méthode HTTP',
  'page.route.form.itemLabel.scheme': 'Schéma',
  'page.route.form.itemLabel.priority': 'Priorité',
  'page.route.form.itemLabel.redirect': 'Redirection',
  'page.route.select.option.enableHttps': 'Activer HTTPS',
  'page.route.select.option.configCustom': 'Personnalisé',
  'page.route.select.option.forbidden': 'Interdit',
  'page.route.form.itemLabel.redirectCustom': 'Redirection personnalisée',
  'page.route.input.placeholder.redirectCustom': 'Par exemple: /foo/index.html',
  'page.route.select.option.redirect301': '301 (Redirection permanente)',
  'page.route.select.option.redirect302': '302 (Redirection temporaire)',
  'page.route.form.itemLabel.username': "Nom d'utilisateur",
  'page.route.form.itemLabel.password': 'Mot de passe',
  'page.route.form.itemLabel.token': 'Jeton',
  'page.route.form.itemLabel.apikey': 'Clé API',
  'page.route.form.itemExtraMessage.domain':
    'Nom de domaine ou IP, prise en charge du nom de domaine générique, par exemple : *.test.com',
  'page.route.form.itemRulesPatternMessage.domain':
    'Seules les lettres, les chiffres, -,_ et * sont pris en charge, mais * doit être au début.',
  'page.route.form.itemExtraMessage1.path':
    'Chemin de la requête HTTP, par exemple : /foo/index.html, prend en charge le préfixe du chemin de la requête /foo/* ; /* représente tous les chemins',
  'page.route.form.itemRulesPatternMessage.remoteAddrs':
    'Veuillez saisir une adresse IP valide, par exemple : 192.168.1.101, 192.168.1.0/24, ::1, fe80::1, fe80::1/64',
  'page.route.form.itemExtraMessage1.remoteAddrs':
    'Adresse IP client, par exemple : 192.168.1.101, 192.168.1.0/24, ::1, fe80::1, fe80::1/64',
  'page.route.httpAction': 'Action',
  'page.route.httpOverrideOrCreate': 'Remplacer/Créer',
  'page.route.panelSection.title.httpOverrideRequestHeader': "Remplacer l'en-tête de requête HTTP",
  'page.route.status': 'Statut',
  'page.route.groupName': 'Nom du groupe',
  'page.route.offline': 'Hors ligne',
  'page.route.publish': 'Publier',
  'page.route.published': 'Publié',
  'page.route.unpublished': 'Non publié',
  'page.route.select.option.inputManually': 'Entrée manuelle',
  'page.route.select.option.methodRewriteNone': 'Ne pas modifier',
  'page.route.form.itemLabel.domainNameOrIp': 'Nom de domaine/IP',
  'page.route.form.itemExtraMessage.domainNameOrIp':
    "Lors de l'utilisation du nom de domaine, il analysera le fichier local: /etc/resolv.conf par défaut",
  'page.route.portNumber': 'Numéro de port',
  'page.route.weight': 'Poids',
  'page.route.radio.staySame': 'Rester identique',
  'page.route.form.itemLabel.newPath': 'Nouveau chemin',
  'page.route.form.itemLabel.newHost': 'Nouveau hôte',
  'page.route.form.itemLabel.regex': 'Expression régulière',
  'page.route.form.itemLabel.template': 'Modèle',
  'page.route.form.itemLabel.URIRewriteType': "Remplacement de l'URI",
  'page.route.form.itemLabel.hostRewriteType': "Remplacement de l'hôte",
  'page.route.form.itemLabel.methodRewrite': 'Remplacement de la méthode',
  'page.route.form.itemLabel.redirectURI': 'URI de redirection',
  'page.route.input.placeholder.newPath': 'Par exemple : /foo/bar/index.html',
  'page.route.steps.stepTitle.defineApiRequest': 'Définir la requête API',
  'page.route.steps.stepTitle.defineApiBackendServe': 'Définir le serveur backend API',
  'page.route.popconfirm.title.offline': 'Êtes-vous sûr de mettre cette route hors ligne?',
  'page.route.radio.static': 'Statique',
  'page.route.radio.regex': 'Regex',
  'page.route.form.itemLabel.from': 'De',
  'page.route.form.itemHelp.status':
    'Indique si une route peut être utilisée après sa création, la valeur par défaut est false.',
  'page.route.host': 'Hôte',
  'page.route.path': 'Chemin',
  'page.route.remoteAddrs': 'Adresses distantes',
  'page.route.PanelSection.title.defineRequestParams': 'Définir les paramètres de requête',
  'page.route.PanelSection.title.responseResult': 'Résultat de la réponse',
  'page.route.debug.showResultAfterSendRequest': "Afficher le résultat après l'envoi de la requête",
  'page.route.TabPane.queryParams': 'Paramètres de requête',
  'page.route.TabPane.bodyParams': 'Paramètres du corps',
  'page.route.TabPane.headerParams': "Paramètres d'en-tête",
  'page.route.TabPane.authentication': 'Authentification',
  'page.route.TabPane.response': 'Réponse',
  'page.route.TabPane.header': 'En-tête de réponse',
  'page.route.debugWithoutAuth': "Cette requête n'utilise aucune autorisation.",
  'page.route.button.exportOpenApi': 'Exporter OpenAPI',
  'page.route.exportRoutesTips': 'Veuillez choisir le type de fichier que vous souhaitez exporter',
  'page.route.button.importOpenApi': 'Importer OpenAPI',
  'page.route.button.selectFile': 'Veuillez sélectionner le fichier',
  'page.route.list': 'Routes',
  'page.route.panelSection.title.requestOverride': 'Remplacement de la requête',
  'page.route.form.itemLabel.headerRewrite': "Remplacement de l'en-tête",
  'page.route.tooltip.pluginOrchOnlySupportChrome':
    "L'orchestration de plugin prend en charge uniquement Chrome.",
  'page.route.tooltip.pluginOrchWithoutProxyRewrite':
    "Le mode d'orchestration de plugin ne peut pas être utilisé lorsque le remplacement de requête est configuré à l'étape 1.",
  'page.route.tooltip.pluginOrchWithoutRedirect':
    "Le mode d'orchestration de plugin ne peut pas être utilisé lorsque la redirection à l'étape 1 est activée pour activer HTTPS.",
  'page.route.tabs.normalMode': 'Normal',
  'page.route.tabs.orchestration': 'Orchestration',
  'page.route.list.description':
    "La route est le point d'entrée d'une requête, qui définit les règles de correspondance entre une requête client et un service. Une route peut être associée à un service (Service), à un amont (Upstream), un service pouvant correspondre à un ensemble de routes et une route pouvant correspondre à un objet amont (un ensemble de nœuds de service backend), de sorte que chaque requête correspondant à une route sera relayée par la passerelle vers le service amont lié à la route.",
  'page.route.configuration.name.rules.required.description': 'Veuillez saisir le nom de la route',
  'page.route.configuration.name.placeholder': 'Veuillez saisir le nom de la route',
  'page.route.configuration.desc.tooltip': 'Veuillez saisir la description de la route',
  'page.route.configuration.publish.tooltip':
    'Utilisé pour contrôler si une route est publiée sur la passerelle immédiatement après sa création',
  'page.route.configuration.version.placeholder':
    'Veuillez saisir le numéro de version de la route',
  'page.route.configuration.version.tooltip': 'Numéro de version de la route, par exemple V1',
  'page.route.configuration.normal-labels.tooltip':
    'Ajoutez des étiquettes personnalisées aux routes pouvant être utilisées pour le regroupement des routes.',
  'page.route.configuration.path.rules.required.description':
    'Veuillez saisir un chemin de requête HTTP valide',
  'page.route.configuration.path.placeholder': 'Veuillez saisir le chemin de requête HTTP',
  'page.route.configuration.remote_addrs.placeholder': "Veuillez saisir l'adresse client",
  'page.route.configuration.host.placeholder': "Veuillez saisir le nom d'hôte de la requête HTTP",
  'page.route.service.none': 'Aucun',
  'page.route.rule.create': 'Créer une règle',
  'page.route.rule.edit': 'Configurer une règle',
  'page.route.advanced-match.operator.sample.IN':
    'Veuillez entrer un tableau, par exemple ["1", "2"]',
  'page.route.advanced-match.operator.sample.~~':
    'Veuillez entrer une expression régulière, par exemple [a-z]+',
  'page.route.fields.service_id.invalid':
    'Veuillez vérifier la configuration de liaison du service',
  'page.route.fields.service_id.without-upstream':
    "Si vous ne liez pas le service, vous devez définir l'amont (Étape 2)",
  'page.route.advanced-match.tooltip':
    'Il prend en charge la correspondance de route via les en-têtes de requête, les paramètres de requête et les cookies, et peut être appliqué à des scénarios tels que la publication en niveaux de gris et les tests bleu-vert.',
  'page.route.advanced-match.message': 'Conseils',
  'page.route.advanced-match.tips.requestParameter': 'Paramètre de requête : Requête URL',
  'page.route.advanced-match.tips.postRequestParameter':
    'Paramètre de requête POST : Prise en charge du formulaire x-www-form-urlencoded',
  'page.route.advanced-match.tips.builtinParameter':
    'Paramètre intégré : Paramètres internes de Nginx',
  'page.route.fields.custom.redirectOption.tooltip': 'Cela concerne le plugin de redirection',
  'page.route.fields.service_id.tooltip':
    "Associez l'objet Service pour réutiliser leur configuration.",
  'page.route.fields.vars.invalid':
    'Veuillez vérifier la configuration de la condition de correspondance avancée',
  'page.route.fields.vars.in.invalid':
    "Lors de l'utilisation de l'opérateur IN, entrez les valeurs des paramètres au format tableau.",
  'page.route.data_loader.import': 'Importer',
  'page.route.data_loader.import_panel': 'Importer des données',
  'page.route.data_loader.types.openapi3': 'OpenAPI 3',
  'page.route.data_loader.types.openapi_legacy': 'OpenAPI 3 Legacy',
  'page.route.data_loader.labels.loader_type': 'Type de chargeur de données',
  'page.route.data_loader.labels.task_name': 'Nom de la tâche',
  'page.route.data_loader.labels.upload': 'Télécharger',
  'page.route.data_loader.labels.openapi3_merge_method': 'Fusionner les méthodes HTTP',
  'page.route.data_loader.tips.select_type': 'Veuillez sélectionner le chargeur de données',
  'page.route.data_loader.tips.input_task_name': "Veuillez saisir le nom de la tâche d'importation",
  'page.route.data_loader.tips.click_upload': 'Cliquez pour télécharger',
  'page.route.data_loader.tips.openapi3_merge_method':
    'Permet de fusionner plusieurs méthodes HTTP dans le chemin OpenAPI en une seule route. Lorsque vous avez plusieurs méthodes HTTP dans votre chemin avec une configuration différente (par exemple, securitySchema), vous pouvez désactiver cette option pour les générer en plusieurs routes.',
};

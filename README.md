# 🔥 Ultimate TCP SYN Flooder with Proxy Support

![TCP-SYN-Flooder](https://www.imperva.com/learn/wp-content/uploads/sites/13/2019/01/syn-flood.jpg)
🚀 **Un puissant script Golang pour mener des attaques TCP SYN Flood tout en utilisant des proxies HTTP !**

---

## ⚡ Fonctionnalités Principales

✅ **Attaque TCP SYN Flood** sur une cible spécifiée 🎯  
✅ **Utilisation automatique de proxies HTTP** pour masquer l'origine 📡  
✅ **Génération d'IP aléatoires** pour éviter la détection 🚀  
✅ **Gestion multi-threadée** pour maximiser l'efficacité 💥  
✅ **Contrôle du débit** via la limitation des PPS 📊  
✅ **Facile à utiliser** avec des arguments en ligne de commande ⌨️  

---

## 🛠️ Installation & Utilisation

### 📌 Prérequis
- **Go** (Dernière version recommandée) → [Téléchargement](https://go.dev/dl/)
- Connexion Internet pour récupérer les proxies 💻

### ⬇️ Installation
```bash
git clone https://github.com/MedusaSH/TCP-SYN-Flooder.git
cd TCP-SYN
go build -o flooder tcpsyn.go
```

### 🚀 Utilisation
```bash
./flooder <IP cible> <port> <threads> <PPS> <durée>
```
📌 **Exemple :**
```bash
./flooder 192.168.1.1 80 100 5000 60
```
🔹 Envoie une attaque sur **192.168.1.1:80** avec **100 threads**, **5000 PPS**, pendant **60 secondes**.

---

## 🎯 Explication du Fonctionnement

⚙️ **Comment ça marche ?**
1. **Récupération des proxies** à partir d'une API publique.
2. **Génération d'adresses IP et ports aléatoires** pour chaque requête.
3. **Construction manuelle des paquets TCP/IP** avec spoofing d'IP.
4. **Envoi de requêtes à haute fréquence** via des threads multiples.
5. **Résistance aux contre-mesures** en utilisant des IP différentes et des flags TCP modifiés.

---

## ⚠️ Avertissement Légal
> **Ce script est fourni uniquement à des fins éducatives.**
> **L'utilisation non autorisée de ce programme peut être illégale.**
> **L'auteur décline toute responsabilité en cas d'utilisation malveillante.** ⚖️

---

## ❤️ Contribuer
👥 **Envie d'améliorer le script ?** N'hésite pas à **fork**, proposer des **PRs** ou signaler des **issues** !

📫 Contact : [TonEmail@exemple.com](mailto:medusa.cc.pro@gmail.com)

🚀 **Star le repo si ce projet t'a été utile !** ⭐


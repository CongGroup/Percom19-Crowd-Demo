<template>
    <div class="home">
        <submit :account="account" />
    </div>
</template>

<script>
    // @ is an alias to /src
    require("@/assets/css/style.css");
    export default {
        name: "home",
        components: {
            submit: () => import("@/components/submit.vue")
        },
        data: function () {
            return {

            }
        },
        computed: {
            account: function () {
                return this.$store.state.account
            }
        },
        beforeCreate: function () {
            import('../assets/js/tx.js').then(util => {
                let account;
                if(this.$cookies.isKey("account")) {
                    console.log("stored in cookie");
                    account = this.$cookies.get("account");
                } else {
                    console.log("create account");
                    account = util.createAccount();
                    this.$cookies.set('account',account,'7d');
                }
                this.$store.commit('setAccount', account);
                account.privateKey = account.privateKey.slice(2);
                //test
                console.log("address: "+account.address);
                console.log("private key: " + account.privateKey);
            })
        }
    };
</script>

<template>
    <div>
        <manager :serviceProviderAccount="serviceProviderAccount" :dataConsumerAccount="dataConsumerAccount"></manager>
    </div>
</template>

<script>
    // @ is an alias to /src
    require("@/assets/css/style.css");
    export default {
        name: "admin",
        components: {
            manager: () => import("../components/manager.vue")
        },
        data: function () {
            return {
                serviceProviderAccount: undefined,
                dataConsumerAccount: undefined,
            }
        },
        computed: {
            account: function () {
                return this.$store.state.account
            }
        },
        created: function() {
            this.serviceProviderAccount = {
                address: "0x88d6F234b6327b7a9d9d23014EFd6E05b6Dc8943",
                privateKey: "65CF458E20E12991BA6AEBC210B8293AD5B3495F585CECC6D22DE77933DE2B0F"
            };
            this.dataConsumerAccount = {
                address: "0x3c62aa7913bc303ee4b9c07df87b556b6770e3fc",
                privateKey: "e27cb51d1eb94ad42b8f196e341e082042639677df43fd7d1440c07b40e2a065"
            };
        }
    };
</script>

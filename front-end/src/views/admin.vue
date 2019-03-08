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
                address: "0x5eb76be14a2408508f64a2d35acf3b4d19198d45",
                privateKey: "275b9b0eb20cbf877c96ab88cbe11b7609229162c15f0f2d206d85492be94f1b"
            };
        }
    };
</script>


techPayControllers.controller('HomeCtrl',  function($scope, $location, $cookieStore, cantineiro) {
    $scope.searchAluno = function () {

        if( $scope.search.length >0 ){
            $scope.usuarios=[]
            for( var i=0;i<$scope.allUsers.length;i++){
                if( $scope.allUsers[i].Nome.toUpperCase().indexOf( $scope.search.toUpperCase() ) != -1 ){
                    $scope.usuarios.push($scope.allUsers[i] );
                }
                if( $scope.allUsers[i].RA.toUpperCase().indexOf( $scope.search.toUpperCase() ) != -1 ){
                    $scope.usuarios.push($scope.allUsers[i] );
                }
            } 
        } else {
            $scope.usuarios = $scope.allUsers ;
        }

    }
    $scope.setUser = function () {
        $scope.creditostoadd="";
        for( i=0;i<$scope.usuarios.length;i++ ){
            if( $scope.usuarios[i].id == $scope.selectedUser ){
                $scope.user = $scope.usuarios[i];
            }
        }
    }
    $scope.addCredit = function () {

        // var loading_screen = pleaseWait({
        //     backgroundColor: '#2277ee',
        //     loadingHtml: "<div class='sk-spinner sk-spinner-wave'><div class='sk-rect1'></div><div class='sk-rect2'></div><div class='sk-rect3'></div><div class='sk-rect4'></div><div class='sk-rect5'></div></div>"
        // });

        // alert( JSON.stringify($scope.selectedUser));

        var credentials = {
            token: $cookieStore.get('token'),
            idaluno: parseInt($scope.selectedUser),
            idtipo: 1,
            valor: $scope.creditostoadd
        };

        var success = function (data) {
//            loading_screen.finish();
//        alert( JSON.stringify( data ) ) ;

            switch( data.message ){
                case 'success':
                    alert("Creditos inseridos com sucesso");
                    break;
                default:
                    alert( data.message );
                    break;
            }
        };

        var error = function (v) {
            alert("ERRO " );
            // TODO: apply user notification here..
        };

        cantineiro.insereCreditos(credentials).success(success).error(error);
    };
    $scope.listAlunos = function () {

        // var loading_screen = pleaseWait({
        //     backgroundColor: '#2277ee',
        //     loadingHtml: "<div class='sk-spinner sk-spinner-wave'><div class='sk-rect1'></div><div class='sk-rect2'></div><div class='sk-rect3'></div><div class='sk-rect4'></div><div class='sk-rect5'></div></div>"
        // });

        var credentials = {
            token: $cookieStore.get('token'),
        };

        var success = function (data) {
//            loading_screen.finish();
//        alert( JSON.stringify( data ) ) ;

            switch( data.message ){
                case 'success':
                    $scope.allUsers = data.data.alunos;
                    $scope.usuarios = data.data.alunos;

                    $scope.selectedUser = $scope.usuarios[1].id;

                    break;
                default:
                    break;
            }
        };

        var error = function (v) {
            alert("ERRO " );
            // TODO: apply user notification here..
        };

        cantineiro.listAlunos(credentials).success(success).error(error);
    };
     $scope.currencySymbol='R$';
     $scope.listAlunos();
});


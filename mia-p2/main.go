package main

import (
	"bufio"
	"container/list"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Main function
func main() {
	ListaDiscos := list.New()
	LlenarListaDisco(ListaDiscos)
	//LeerArchivo de Entrada
	var comando string = ""
	scanner := bufio.NewScanner(os.Stdin)
	for comando != "exit" {
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Println("                            MIA PROYECTO2 EC2                                   ")
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Print(">>")
		scanner.Scan()
		comando = scanner.Text()
		if comando != "" {
			LeerTexto(comando, ListaDiscos)
		}

	}
}

func LlenarListaDisco(ListaDiscos *list.List) {
	IdDisco := [26]string{"a", "b", "c", "d", "e", "f", "g", "h", "i",
		"j", "k", "l", "m", "n", "o", "p", "q",
		"r", "s", "t", "u", "v", "w", "x", "y", "z"}
	for i := 0; i < 26; i++ {
		disco := DISCO{}
		copy(disco.Estado[:], "0")
		copy(disco.Id[:], IdDisco[i])
		for j := 0; j < len(disco.Particiones); j++ {
			mount := MOUNT{}
			mount.NombreParticion = ""
			mount.Id = strconv.Itoa(j + 1)
			copy(mount.Estado[:], "0")
			disco.Particiones[j] = mount
		}
		ListaDiscos.PushBack(disco)
	}
}

var global string = ""

//Funcion para leer y reconocer los comandos lleno la lista de comandos
func LeerTexto(dat string, ListaDiscos *list.List) {
	//Leendo la cadena de entrada
	ListaComandos := list.New()
	lineaComando := strings.Split(dat, "\n")
	var c Comando
	for i := 0; i < len(lineaComando); i++ {
		EsComentario := lineaComando[i][0:1]
		if EsComentario != "#" {
			comando := lineaComando[i]
			if strings.Contains(lineaComando[i], "\\*") {
				comando = strings.Replace(lineaComando[i], "\\*", " ", 1) + lineaComando[i+1]
				i = i + 1
			}
			propiedades := strings.Split(string(comando), " ")
			nombreComando := propiedades[0]
			c = Comando{Name: strings.ToLower(nombreComando)}
			propiedadesTemp := make([]Propiedad, len(propiedades)-1)
			for i := 1; i < len(propiedades); i++ {
				if propiedades[i] == "" {
					continue
				} else if propiedades[i] == "-p" {
					propiedadesTemp[i-1] = Propiedad{Name: "-p",
						Val: "-p"}
				} else {
					if strings.Contains(propiedades[i], "=") {
						valor_propiedad_Comando := strings.Split(propiedades[i], "=")
						propiedadesTemp[i-1] = Propiedad{Name: valor_propiedad_Comando[0],
							Val: valor_propiedad_Comando[1]}
					} else {
						propiedadesTemp[i-1] = Propiedad{Name: "-sigue",
							Val: propiedades[i]}
					}
				}
			}
			c.Propiedades = propiedadesTemp
			ListaComandos.PushBack(c)
		}
	}
	RecorrerListaComando(ListaComandos, ListaDiscos)
}

func RecorrerListaComando(ListaComandos *list.List, ListaDiscos *list.List) {
	var ParamValidos bool = true
	var cont = 1
	for element := ListaComandos.Front(); element != nil; element = element.Next() {
		var comandoTemp Comando
		comandoTemp = element.Value.(Comando)
		switch strings.ToLower(comandoTemp.Name) {
		case "mkdisk":
			ParamValidos = EjecutarComandoMKDISK(comandoTemp.Name, comandoTemp.Propiedades, cont)
			cont++
			if ParamValidos == false {
				fmt.Println("Parametros Invalidos")
			}
		case "rmdisk":
			//ParamValidos = EjecutarComandoRMDISK(comandoTemp.Name, comandoTemp.Propiedades)
			if ParamValidos == false {
				fmt.Println("Parametros Invalidos")
			}
		case "fdisk":
			//ParamValidos = EjecutarComandoFDISK(comandoTemp.Name, comandoTemp.Propiedades)
			if ParamValidos == false {
				fmt.Println("Parametros Invalidos")
			}
		case "mount":
			/*if len(comandoTemp.Propiedades) != 0 {
				ParamValidos = EjecutarComandoMount(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
				if ParamValidos == false {
					fmt.Println("Parametros Invalidos")
				}
			} else {
				EjecutarReporteMount(ListaDiscos)
			}*/
		case "exit":
			fmt.Println("")
		case "pause":
			fmt.Println("Presione la tecla enter para Continuar")
			fmt.Scanln()
		case "exec":
			ParamValidos = EjecutarComandoExec(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("Error, parametros no validos")
			}
		case "mkdir":
			//ParamValidos = EjecutarComandoMKDIR(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("Error, parametros no validos")
			}
		case "mkfile":
			//ParamValidos = EjecutarComandoMKFILE(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("Error, parametros no validos")
			}
		case "mkfs":
			//ParamValidos = EjecutarComandoMKFS(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("Error, parametros no validos")
			}
		case "rep":
			//ParamValidos = EjecutarComandoReporte(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("Error, parametros no validos")
			}
		case "rmgrp":
			//ParamValidos = EjecutarComandoMKGRP(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("Error, parametros no validos")
			}
		case "login":
			//ParamValidos, global = EjecutarComandoLogin(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
			if ParamValidos == false {
				fmt.Println("Error, parametros no validos")
			}
		case "logout":
			/*if global == "" {
				fmt.Println("Inicia sesión primero para realizar la acción.")
			} else {
				global = ""
				fmt.Println("Sesión finalizada")
			}*/
		default:
			fmt.Println("Error, comando invalido")
		}
	}
}

/*------------------------------MKDISK-------------------------------------------*/
func EjecutarComandoMKDISK(nombreComando string, propiedadesTemp []Propiedad, cont int) (ParamValidos bool) {
	dt := time.Now()
	mbr1 := MBR{}
	copy(mbr1.MbrFechaCreacion[:], dt.String())
	mbr1.NoIdentificador = int64(rand.Intn(100) + cont)
	fmt.Println("----------------- Ejecutando MKDISK -----------------")
	comandos := "dd if=/dev/zero "
	ParamValidos = true
	pathCompleta := ""
	var propiedades [4]string
	if len(propiedadesTemp) >= 2 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			//Vector temporal de propiedades
			switch strings.ToLower(nombrePropiedad) {
			case "-size":
				propiedades[0] = propiedadTemp.Val
			case "-unit":
				propiedades[2] = strings.ToLower(propiedadTemp.Val)
			case "-path":
				propiedades[3] = propiedadTemp.Val
				arr_path := strings.Split(propiedades[3], "/")

				for i := 0; i < len(arr_path)-1; i++ {
					pathCompleta += arr_path[i] + "/"
				}
				executeComand("mkdir " + pathCompleta)
				comandos += "of=" + propiedades[3]
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		EsComilla := propiedades[3][0:1]
		if EsComilla == "\"" {
			propiedades[3] = propiedades[3][1 : len(propiedades[3])-1]
		}
		tamanioTotal, _ := strconv.ParseInt(propiedades[0], 10, 64)
		if propiedades[2] == "k" {
			comandos += " bs=" + strconv.Itoa((int(tamanioTotal))*1000) + " count=1"
			mbr1.MbrTamanio = ((tamanioTotal) - 1) * 1000
		} else {
			comandos += " bs=" + strconv.Itoa(int(tamanioTotal)) + "MB" + " count=1"
			mbr1.MbrTamanio = tamanioTotal * 1000000
		}
		//Inicializando Particiones
		for i := 0; i < 4; i++ {
			copy(mbr1.Particiones[i].Status_particion[:], "0")
			copy(mbr1.Particiones[i].TipoParticion[:], "")
			copy(mbr1.Particiones[i].TipoAjuste[:], "")
			mbr1.Particiones[i].Inicio_particion = 0
			mbr1.Particiones[i].TamanioTotal = 0
			copy(mbr1.Particiones[i].NombreParticion[:], "")
		}

		executeComand(comandos)

		f, err := os.OpenFile(propiedades[3], os.O_WRONLY, 0755)
		if err != nil {
			log.Fatalln(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatalln(err)
			}
		}()
		f.Seek(0, 0)
		err = binary.Write(f, binary.BigEndian, mbr1)
		if err != nil {
			log.Fatalln(err, propiedades[3])
		}
		fmt.Println("Se ha creado el disco en la carpeta: ", pathCompleta)
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

/*-----------------------EXEC----------------------------------------------*/
func EjecutarComandoExec(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	fmt.Println("----------------- Ejecutando EXEC -----------------")
	ParamValidos = true
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-path":
				fmt.Println(propiedadTemp.Val)
				dat, err := ioutil.ReadFile(propiedadTemp.Val)
				CheckError(err)
				LeerTexto(string(dat), ListaDiscos)
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

func executeComand(comandos string) {
	args := strings.Split(comandos, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.CombinedOutput()
}
func BytesToString(data [1]byte) string {
	return string(data[:])
}
func CheckError(e error) {
	if e != nil {
		fmt.Println("Error - ----------")
		fmt.Println(e)
	}
}

//Estructura para cada Comando y sus Propiedades
type Propiedad struct {
	Name string
	Val  string
}
type Comando struct {
	Name        string
	Propiedades []Propiedad
}

//Estructuras para el Disco y Particiones
type Particion struct {
	Status_particion [1]byte
	TipoParticion    [1]byte
	TipoAjuste       [2]byte
	Inicio_particion int64
	TamanioTotal     int64
	NombreParticion  [15]byte
}

//Struct para el MBR
type MBR struct {
	MbrTamanio       int64
	MbrFechaCreacion [19]byte
	NoIdentificador  int64
	Particiones      [4]Particion
}

//Struct para las particiones Logicas
type EBR struct {
	Status_particion    [1]byte
	TipoAjuste          [2]byte
	Inicio_particion    int64
	Particion_Siguiente int64
	TamanioTotal        int64
	NombreParticion     [15]byte
}

//EStruc de las particiones montadas
type MOUNT struct {
	NombreParticion string
	Id              string
	Estado          [1]byte
	EstadoMKS       [1]byte
}

//Estruct Disco
type DISCO struct {
	NombreDisco string
	Path        string
	Id          [1]byte
	Estado      [1]byte
	Particiones [100]MOUNT
}

//57.51
type Integers struct {
	I1  uint16
	I2  int32
	I3  int64
	DOS byte
}

type SB struct {
	Sb_nombre_hd                          [15]byte
	Sb_arbol_virtual_count                int64
	Sb_detalle_directorio_count           int64
	Sb_inodos_count                       int64
	Sb_bloques_count                      int64
	Sb_arbol_virtual_free                 int64
	Sb_detalle_directorio_free            int64
	Sb_inodos_free                        int64
	Sb_bloques_free                       int64
	Sb_date_creacion                      [19]byte
	Sb_date_ultimo_montaje                [19]byte
	Sb_montajes_count                     int64
	Sb_ap_bitmap_arbol_directorio         int64
	Sb_ap_arbol_directorio                int64
	Sb_ap_bitmap_detalle_directorio       int64
	Sb_ap_detalle_directorio              int64
	Sb_ap_bitmap_tabla_inodo              int64
	Sb_ap_tabla_inodo                     int64
	Sb_ap_bitmap_bloques                  int64
	Sb_ap_bloques                         int64
	Sb_ap_log                             int64
	Sb_size_struct_arbol_directorio       int64
	Sb_size_struct_Detalle_directorio     int64
	Sb_size_struct_inodo                  int64
	Sb_size_struct_bloque                 int64
	Sb_first_free_bit_arbol_directorio    int64
	Sb_first_free_bit_detalle_directoriio int64
	Sb_dirst_free_bit_tabla_inodo         int64
	Sb_first_free_bit_bloques             int64
	Sb_magic_num                          int64
	InicioCopiaSB                         int64
	ConteoAVD                             int64
	ConteoDD                              int64
	ConteoInodo                           int64
	ConteoBloque                          int64
}

//Detalle dde Directorio

type ArregloDD struct {
	Dd_file_nombre            [15]byte
	Dd_file_ap_inodo          int64
	Dd_file_date_creacion     [19]byte
	Dd_file_date_modificacion [19]byte
}
type DD struct {
	Dd_array_files           [5]ArregloDD
	Dd_ap_detalle_directorio int64
	Ocupado                  int8
}

//Cantidad de Inodos
type Inodo struct {
	I_count_inodo             int64
	I_size_archivo            int64
	I_count_bloques_asignados int64
	I_array_bloques           [4]int64
	I_ao_indirecto            int64
	I_id_proper               int64
}

//Bloque
type Bloque struct {
	Db_data [25]byte
}

//
func BytesNombreParticion(data [15]byte) string {
	return string(data[:])
}
func ConvertData(data [25]byte) string {
	return string(data[:])
}

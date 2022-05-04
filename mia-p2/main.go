package main

import (
	"bufio"
	"bytes"
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
	"unsafe"
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
			ParamValidos = EjecutarComandoFDISK(comandoTemp.Name, comandoTemp.Propiedades)
			if ParamValidos == false {
				fmt.Println("Parametros Invalidos")
			}
		case "mount":
			if len(comandoTemp.Propiedades) != 0 {
				ParamValidos = EjecutarComandoMount(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
				if ParamValidos == false {
					fmt.Println("Parametros Invalidos")
				} else {
					EjecutarReporteMount(ListaDiscos)
				}
			}
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
			ParamValidos = EjecutarComandoReporte(comandoTemp.Name, comandoTemp.Propiedades, ListaDiscos)
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

/*----------------------------FDISK---------------------------------------*/
func EjecutarComandoFDISK(nombreComando string, propiedadesTemp []Propiedad) (ParamValidos bool) {
	fmt.Println("----------------- Ejecutando FDISK -----------------")
	ParamValidos = true
	mbr := MBR{}
	particion := Particion{}
	var startPart int64 = int64(unsafe.Sizeof(mbr))
	var propiedades [8]string
	if len(propiedadesTemp) >= 2 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-size":
				propiedades[0] = propiedadTemp.Val
			case "-fit":
				propiedades[1] = propiedadTemp.Val
			case "-unit":
				propiedades[2] = propiedadTemp.Val
			case "-path":
				propiedades[3] = propiedadTemp.Val
			case "-type":
				propiedades[4] = propiedadTemp.Val
			case "-delete":
				propiedades[5] = propiedadTemp.Val
			case "-name":
				propiedades[6] = propiedadTemp.Val
				fmt.Println(propiedades[6])
			case "-add":
				propiedades[7] = propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		EsComilla := propiedades[3][0:1]
		if EsComilla == "\"" {
			propiedades[3] = propiedades[3][1 : len(propiedades[3])-1]
		}
		//Tamanio Particion
		var TamanioTotalParticion int64 = 0
		if strings.ToLower(propiedades[2]) == "b" {
			TamanioParticion, _ := strconv.ParseInt(propiedades[0], 10, 64)
			TamanioTotalParticion = TamanioParticion
		} else if strings.ToLower(propiedades[2]) == "k" {
			TamanioParticion, _ := strconv.ParseInt(propiedades[0], 10, 64)
			TamanioTotalParticion = TamanioParticion * 1000
		} else if strings.ToLower(propiedades[2]) == "m" {
			TamanioParticion, _ := strconv.ParseInt(propiedades[0], 10, 64)
			TamanioTotalParticion = TamanioParticion * 1000000
		} else {
			TamanioParticion, _ := strconv.ParseInt(propiedades[0], 10, 64)
			TamanioTotalParticion = TamanioParticion * 1000
		}
		if propiedades[5] != "" {
			EliminarParticion(propiedades[3], propiedades[6], propiedades[5])
			return
		}
		//Obtener el MBR
		switch strings.ToLower(propiedades[4]) {
		case "p":
			var Particiones [4]Particion
			f, err := os.OpenFile(propiedades[3], os.O_RDWR, 0755)
			if err != nil {
				fmt.Println("No existe la ruta" + propiedades[3])
				return false
			}
			defer f.Close()
			f.Seek(0, 0)
			err = binary.Read(f, binary.BigEndian, &mbr)
			Particiones = mbr.Particiones
			if err != nil {
				fmt.Println("No existe el archivo en la ruta")
			}

			if HayEspacio(TamanioTotalParticion, mbr.MbrTamanio) {
				return false
			} //Verificar si ya hay particiones
			if BytesToString(Particiones[0].Status_particion) == "1" {
				for i := 0; i < 4; i++ {
					//Posicion en bytes del partstar de la n particion
					startPart += Particiones[i].TamanioTotal
					if BytesToString(Particiones[i].Status_particion) == "0" {
						//fmt.Println(startPart)
						break
					}
				}
			}
			if HayEspacio(startPart+TamanioTotalParticion, mbr.MbrTamanio) {
				return false
			}
			//dando valores a la particion
			copy(particion.Status_particion[:], "1")
			copy(particion.TipoParticion[:], propiedades[4])
			copy(particion.TipoAjuste[:], propiedades[1])
			particion.Inicio_particion = startPart
			particion.TamanioTotal = TamanioTotalParticion
			copy(particion.NombreParticion[:], propiedades[6])
			//Particion creada
			for i := 0; i < 4; i++ {
				if BytesToString(Particiones[i].Status_particion) == "0" {
					Particiones[i] = particion
					break
				}
			}
			f.Seek(0, 0)
			mbr.Particiones = Particiones
			err = binary.Write(f, binary.BigEndian, mbr)
			ReadFile(propiedades[3])
		case "l":
			fmt.Println("Particion Logica")
			if !HayExtendida(propiedades[3]) {
				fmt.Println("No existe una particion Extendida")
				return false
			}
			ebr := EBR{}
			copy(ebr.Status_particion[:], "1")
			copy(ebr.TipoAjuste[:], propiedades[1])
			ebr.Inicio_particion = startPart
			ebr.Particion_Siguiente = 0
			ebr.TamanioTotal = TamanioTotalParticion
			copy(ebr.NombreParticion[:], propiedades[6])
			//Obteniendo el byte donde empezara la particion Logica
			InicioParticionLogica(propiedades[3], ebr)
		case "e":
			//Particiones Extendidas
			var Particiones [4]Particion
			f, err := os.OpenFile(propiedades[3], os.O_RDWR, 0755)
			if err != nil {
				fmt.Println("No existe la ruta" + propiedades[3])
				return false
			}
			defer f.Close()
			f.Seek(0, 0)
			err = binary.Read(f, binary.BigEndian, &mbr)
			Particiones = mbr.Particiones
			if err != nil {
				fmt.Println("No existe el archivo en la ruta")
			}
			//El mbr ya se a leido,2.Verificar si existe espacion disponible o que no lo rebase
			if HayEspacio(TamanioTotalParticion, mbr.MbrTamanio) {
				return false
			} //Verificar si ya hay particiones
			if BytesToString(Particiones[0].Status_particion) == "1" {
				for i := 0; i < 4; i++ {
					//Posicion en bytes del partstar de la n particion
					startPart += Particiones[i].TamanioTotal
					if BytesToString(Particiones[i].Status_particion) == "0" {

						break
					}
				}
			}
			if HayEspacio(startPart+TamanioTotalParticion, mbr.MbrTamanio) {
				return false
			}
			//dando valores a la particion
			copy(particion.Status_particion[:], "1")
			copy(particion.TipoParticion[:], propiedades[4])
			copy(particion.TipoAjuste[:], propiedades[1])
			particion.Inicio_particion = startPart
			particion.TamanioTotal = TamanioTotalParticion
			copy(particion.NombreParticion[:], propiedades[6])
			//Particion creada
			for i := 0; i < 4; i++ {
				if BytesToString(Particiones[i].Status_particion) == "0" {
					Particiones[i] = particion
					break
				}
			}
			f.Seek(0, 0)
			mbr.Particiones = Particiones
			err = binary.Write(f, binary.BigEndian, mbr)
			ReadFile(propiedades[3])
			ebr := EBR{}
			copy(ebr.Status_particion[:], "1")
			copy(ebr.TipoAjuste[:], propiedades[1])
			ebr.Inicio_particion = startPart
			ebr.Particion_Siguiente = -1
			ebr.TamanioTotal = TamanioTotalParticion
			copy(ebr.NombreParticion[:], propiedades[6])
			f.Seek(ebr.Inicio_particion, 0)
			err = binary.Write(f, binary.BigEndian, ebr)

			fmt.Println("Extendida", "Leendo EBR")
		default:
			fmt.Println("Ocurrio un error")
		}
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}
func EscribirParticionLogica(path string, ebr EBR, inicioParticionLogica int64, inicioParticionExtendida int64) bool {
	ebr.Inicio_particion = inicioParticionLogica
	ebr.Particion_Siguiente = inicioParticionLogica + ebr.TamanioTotal
	return true
}

func EliminarParticion(path string, name string, typeDelete string) bool {
	var name2 [15]byte
	Encontrada := false
	copy(name2[:], name)
	f, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	//Posiciono al inicio el Puntero
	f.Seek(0, 0)
	//Leo el mbr
	err = binary.Read(f, binary.BigEndian, &mbr)
	Particiones := mbr.Particiones
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" && BytesNombreParticion(Particiones[i].NombreParticion) == BytesNombreParticion(name2) {
			fmt.Println("Es una Extendida")
			Encontrada = true
		} else if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "p" && BytesNombreParticion(Particiones[i].NombreParticion) == BytesNombreParticion(name2) {
			var partTemp = Particion{}
			copy(partTemp.Status_particion[:], "0")
			copy(partTemp.TipoParticion[:], "")
			copy(partTemp.TipoAjuste[:], "")
			partTemp.Inicio_particion = 0
			partTemp.TamanioTotal = 0
			copy(partTemp.NombreParticion[:], "")
			Particiones[i] = partTemp
			mbr.Particiones = Particiones
			f.Seek(0, 0)
			err = binary.Write(f, binary.BigEndian, &mbr)
			fmt.Println("Particion Primaria Eliminada")
			ReadFile(path)
			Encontrada = true
		}
	}
	if Encontrada == false {
		for i := 0; i < 4; i++ {
			if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" {
				var InicioExtendida int64 = Particiones[i].Inicio_particion
				f.Seek(InicioExtendida, 0)
				ebrAnterior := EBR{}
				ebr := EBR{}
				ebrAnterior = ebr
				err = binary.Read(f, binary.BigEndian, &ebr)
				if ebr.Particion_Siguiente == -1 {
					fmt.Println("No Hay particiones Logicas")
				} else {
					f.Seek(InicioExtendida, 0)
					err = binary.Read(f, binary.BigEndian, &ebr)
					for {
						if BytesNombreParticion(ebr.NombreParticion) == BytesNombreParticion(name2) {
							fmt.Println("Particion Logica Encontrada")
							if strings.ToLower(typeDelete) == "fast" {
								ebrAnterior.Particion_Siguiente = ebr.Particion_Siguiente
								f.Seek(ebrAnterior.Inicio_particion, 0)
								err = binary.Write(f, binary.BigEndian, ebrAnterior)

							} else if strings.ToLower(typeDelete) == "full" {
								ebrAnterior.Particion_Siguiente = ebr.Particion_Siguiente
								f.Seek(ebrAnterior.Inicio_particion, 0)
								err = binary.Write(f, binary.BigEndian, ebrAnterior)
							}
							Encontrada = true
						}
						if ebr.Particion_Siguiente == -1 {
							break
						} else {
							f.Seek(ebr.Particion_Siguiente, 0)
							ebrAnterior = ebr
							err = binary.Read(f, binary.BigEndian, &ebr)
						}
					}
				}
			}
		}
	}
	if Encontrada == false {
		fmt.Println("Error No se encontro la particon")
	}
	return false
}
func InicioParticionLogica(path string, ebr2 EBR) bool {
	f, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	Particiones := mbr.Particiones
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" {
			var InicioExtendida int64 = Particiones[i].Inicio_particion
			f.Seek(InicioExtendida, 0)
			ebr := EBR{}
			err = binary.Read(f, binary.BigEndian, &ebr)
			if ebr.Particion_Siguiente == -1 {
				ebr.Particion_Siguiente = ebr.Inicio_particion + int64(unsafe.Sizeof(ebr)) + ebr2.TamanioTotal
				f.Seek(InicioExtendida, 0)
				err = binary.Write(f, binary.BigEndian, ebr)
				ebr2.Inicio_particion = ebr.Particion_Siguiente
				ebr2.Particion_Siguiente = -1
				f.Seek(ebr2.Inicio_particion, 0)
				err = binary.Write(f, binary.BigEndian, ebr2)

				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				fmt.Println(ebr.Inicio_particion)
				fmt.Println(ebr.Particion_Siguiente)
				return false
			} else {

				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				for {
					if ebr.Particion_Siguiente == -1 {
						//fmt.Println("Es la ultima logica")
						ebr.Particion_Siguiente = ebr.Inicio_particion + int64(unsafe.Sizeof(ebr)) + ebr2.TamanioTotal
						f.Seek(ebr.Inicio_particion, 0)
						err = binary.Write(f, binary.BigEndian, ebr)
						ebr2.Inicio_particion = ebr.Particion_Siguiente
						ebr2.Particion_Siguiente = -1
						f.Seek(ebr2.Inicio_particion, 0)
						err = binary.Write(f, binary.BigEndian, ebr2)
						fmt.Printf("NombreLogica: %s\n", ebr2.NombreParticion)
						break
					} else {
						f.Seek(ebr.Particion_Siguiente, 0)
						err = binary.Read(f, binary.BigEndian, &ebr)
						fmt.Printf("NombreLogica: %s\n", ebr.NombreParticion)
					}

				}
				return false
			}
		}
	}
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}

	return false
}
func HayExtendida(path string) bool {
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	Particiones := mbr.Particiones
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" {
			return true
		}
	}
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}

	return false
}
func ReadFileEBR(path string) (funciona bool) {
	//fmt.Println("****************Leendo EL EBR")
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	Particiones := mbr.Particiones
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	for i := 0; i < 4; i++ {
		if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" {
			var InicioExtendida int64 = Particiones[i].Inicio_particion
			f.Seek(InicioExtendida, 0)
			ebr := EBR{}
			err = binary.Read(f, binary.BigEndian, &ebr)
			if ebr.Particion_Siguiente == -1 {
				fmt.Println("No Hay particiones Logicas")
			} else {
				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				for {
					if ebr.Particion_Siguiente == -1 {
						break
					} else {
						f.Seek(ebr.Particion_Siguiente, 0)
						err = binary.Read(f, binary.BigEndian, &ebr)
					}
					fmt.Printf("NombreLogica: %s\n", ebr.NombreParticion)

				}
			}
		}
	}
	return true
}
func ReadFile(path string) (funciona bool) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	return true
}
func HayEspacio(TamanioTotalParticion int64, tamanioDisco int64) bool {
	if ((TamanioTotalParticion) > tamanioDisco) || (TamanioTotalParticion < 0) {
		fmt.Println("Error Fdisk, El Tamanio de la particion es mayor a el tamanio del disco o el tamanio es incorrecto")
		return true
	}
	return false
}

/*------------------------------MOUNT----------------------------------------*/
func EjecutarComandoMount(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	var propiedades [2]string
	var nombre [15]byte
	ParamValidos = true
	if len(propiedadesTemp) >= 2 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-name":
				propiedades[0] = propiedadTemp.Val
				copy(nombre[:], propiedades[0])
			case "-path":
				propiedades[1] = propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el Comando")
			}
		}
		//Empezar a montar las Particiones
		EjecutarComando(propiedades[1], nombre, ListaDiscos)
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}
func EjecutarReporteMount(ListaDiscos *list.List) {
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco DISCO
		disco = element.Value.(DISCO)
		if disco.NombreDisco != "" {
			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				if mountTemp.NombreParticion != "" {
					fmt.Println("id=", mountTemp.Id, " ruta=", disco.Path, "nombre_particion=", mountTemp.NombreParticion)
				}
			}
		}
	}
}
func IdValido(id string, ListaDiscos *list.List) bool {
	esta := false
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco DISCO
		disco = element.Value.(DISCO)
		if disco.NombreDisco != "" {
			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				if mountTemp.NombreParticion != "" {
					if mountTemp.Id == id {
						return true
					}
				}
			}
		}
	}
	return esta
}
func EjecutarComando(path string, NombreParticion [15]byte, ListaDiscos *list.List) bool {
	var encontrada = false
	lineaComando := strings.Split(path, "/")
	nombreDisco := lineaComando[len(lineaComando)-1]
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + path)
		return false
	}
	defer f.Close()
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	Particiones := mbr.Particiones
	for i := 0; i < 4; i++ {
		if string(Particiones[i].NombreParticion[:]) == string(NombreParticion[:]) {
			encontrada = true
			if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" {
				fmt.Println("Error no se puede Montar una particion Extendida")
			} else {
				ParticionMontar(ListaDiscos, string(NombreParticion[:]), string(nombreDisco), path)
			}
		}
		if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" {
			ebr := EBR{}
			f.Seek(Particiones[i].Inicio_particion, 0)
			err = binary.Read(f, binary.BigEndian, &ebr)
			for {
				if ebr.Particion_Siguiente == -1 {
					break
				} else {
					f.Seek(ebr.Particion_Siguiente, 0)
					err = binary.Read(f, binary.BigEndian, &ebr)
				}
				var nombre string = string(ebr.NombreParticion[:])
				var nombre2 string = string(NombreParticion[:])
				if nombre == nombre2 {
					encontrada = true
					//Montar Particion
					ParticionMontar(ListaDiscos, string(NombreParticion[:]), string(nombreDisco), path)
				}
			}
		}
	}
	if encontrada == false {
		fmt.Println("Error no se encontro la particion")
	}
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	}
	return true
}

func ParticionMontar(ListaDiscos *list.List, nombreParticion string, nombreDisco string, path string) {

	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco DISCO
		disco = element.Value.(DISCO)
		if BytesToString(disco.Estado) == "0" && !ExisteDisco(ListaDiscos, nombreDisco) {
			disco.NombreDisco = nombreDisco
			disco.Path = path
			copy(disco.Estado[:], "1")
			//#id->vda1
			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				if BytesToString(mountTemp.Estado) == "0" {
					mountTemp.Id = "14" + BytesToString(disco.Id) + mountTemp.Id
					mountTemp.NombreParticion = nombreParticion
					copy(mountTemp.Estado[:], "1")
					copy(mountTemp.EstadoMKS[:], "0")
					disco.Particiones[i] = mountTemp
					break
				} else if BytesToString(mountTemp.Estado) == "1" && mountTemp.NombreParticion == nombreParticion {
					break
				}
			}
			element.Value = disco
			break
		} else if BytesToString(disco.Estado) == "1" && ExisteDisco(ListaDiscos, nombreDisco) && nombreDisco == disco.NombreDisco {

			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				if BytesToString(mountTemp.Estado) == "0" {
					mountTemp.Id = "14" + BytesToString(disco.Id) + mountTemp.Id
					mountTemp.NombreParticion = nombreParticion
					copy(mountTemp.Estado[:], "1")
					copy(mountTemp.EstadoMKS[:], "0")
					disco.Particiones[i] = mountTemp
					break
				} else if BytesToString(mountTemp.Estado) == "1" && mountTemp.NombreParticion == nombreParticion {
					//fmt.Println("La Particion ya esta montada")
					break
				}
			}
			element.Value = disco
			break
		}
	}
}
func ExisteDisco(ListaDiscos *list.List, nombreDisco string) bool {
	Existe := false
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco DISCO
		disco = element.Value.(DISCO)
		if disco.NombreDisco == nombreDisco {
			return true
		} else {
			Existe = false
		}
	}
	return Existe
}

/*--------------------------REPORTES----------------------------------------*/
func EjecutarComandoReporte(nombreComando string, propiedadesTemp []Propiedad, ListaDiscos *list.List) (ParamValidos bool) {
	ParamValidos = true
	var propiedades [4]string
	if len(propiedadesTemp) >= 1 {
		//Recorrer la lista de propiedades
		for i := 0; i < len(propiedadesTemp); i++ {
			var propiedadTemp = propiedadesTemp[i]
			var nombrePropiedad string = propiedadTemp.Name
			switch strings.ToLower(nombrePropiedad) {
			case "-id":
				propiedades[0] = propiedadTemp.Val
			case "-path":
				propiedades[1] = propiedadTemp.Val
			case "-name":
				propiedades[2] = propiedadTemp.Val
			case "-ruta":
				propiedades[3] = propiedadTemp.Val
			case "-sigue":
				propiedades[1] += propiedadTemp.Val
			default:
				fmt.Println("Error al Ejecutar el Comando", nombrePropiedad)
			}
		}
		EsComilla := propiedades[1][0:1]
		if EsComilla == "\"" {
			if propiedades[3] != "" {
				propiedades[3] = propiedades[3][1 : len(propiedades[3])-1]
			}
			propiedades[1] = propiedades[1][1 : len(propiedades[1])-1]
		}
		carpetas_Graficar := strings.Split(propiedades[1], "/")
		var comando = ""
		for i := 1; i < len(carpetas_Graficar)-1; i++ {
			comando += carpetas_Graficar[i] + "/"
		}
		fmt.Println(comando)
		executeComand("mkdir " + comando[0:len(comando)-1])
		switch strings.ToLower(propiedades[2]) {
		case "mbr":
			GraficarDisco(propiedades[0], ListaDiscos, propiedades[1])
		case "tree":
			//GraficarTreeFull(propiedades[0], propiedades[1], propiedades[3], ListaDiscos)
		default:
			fmt.Println("Reporte incorrecto.")

		}
		return ParamValidos
	} else {
		ParamValidos = false
		return ParamValidos
	}
}

func GraficarDisco(idParticion string, ListaDiscos *list.List, path string) bool {
	var NombreParticion [15]byte
	var buffer bytes.Buffer
	buffer.WriteString("digraph G{\ntbl [\nshape=box\nlabel=<\n<table border='0' cellborder='0' width='100' height=\"30\">\n<tr>")
	pathDisco, _, _ := RecorrerListaDisco(idParticion, ListaDiscos)
	fmt.Println(pathDisco)
	f, err := os.OpenFile(pathDisco, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("No existe la ruta" + pathDisco)
		return false
	}
	defer f.Close()
	PorcentajeUtilizao := 0.0
	var EspacioUtilizado int64 = 0
	mbr := MBR{}
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbr)
	TamanioDisco := mbr.MbrTamanio
	Particiones := mbr.Particiones
	buffer.WriteString("<td height='30' width='75'> MBR </td>")
	for i := 0; i < 4; i++ {
		if convertName(Particiones[i].NombreParticion[:]) != convertName(NombreParticion[:]) && strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "p" {
			PorcentajeUtilizao = (float64(Particiones[i].TamanioTotal) / float64(TamanioDisco)) * 100
			buffer.WriteString("<td height='30' width='75.0'>PRIMARIA <br/>" + convertName(Particiones[i].NombreParticion[:]) + " <br/> Ocupado: " + strconv.Itoa(int(PorcentajeUtilizao)) + "%</td>")
			EspacioUtilizado += Particiones[i].TamanioTotal
		} else if convertName(Particiones[i].Status_particion[:]) == "0" {
			buffer.WriteString("<td height='30' width='75.0'>Libre</td>")
		}
		if strings.ToLower(BytesToString(Particiones[i].TipoParticion)) == "e" {
			EspacioUtilizado += Particiones[i].TamanioTotal
			PorcentajeUtilizao = (float64(Particiones[i].TamanioTotal) / float64(TamanioDisco)) * 100
			buffer.WriteString("<td  height='30' width='15.0'>\n")
			buffer.WriteString("<table border='5'  height='30' WIDTH='15.0' cellborder='1'>\n")
			buffer.WriteString(" <tr>  <td height='60' colspan='100%'>EXTENDIDA <br/>" + convertName(Particiones[i].NombreParticion[:]) + " <br/> Ocupado:" + strconv.Itoa(int(PorcentajeUtilizao)) + "%</td>  </tr>\n<tr>")
			var InicioExtendida int64 = Particiones[i].Inicio_particion
			f.Seek(InicioExtendida, 0)
			ebr := EBR{}
			err = binary.Read(f, binary.BigEndian, &ebr)
			if ebr.Particion_Siguiente == -1 {
				fmt.Println("No Hay particiones Logicas")
			} else {
				var EspacioUtilizado int64 = 0
				cont := 0
				f.Seek(InicioExtendida, 0)
				err = binary.Read(f, binary.BigEndian, &ebr)
				for {
					if ebr.Particion_Siguiente == -1 {
						break
					} else {
						f.Seek(ebr.Particion_Siguiente, 0)
						err = binary.Read(f, binary.BigEndian, &ebr)
						EspacioUtilizado += ebr.TamanioTotal
						PorcentajeUtilizao = (float64(ebr.TamanioTotal) / float64(Particiones[i].TamanioTotal)) * 100
						buffer.WriteString("<td height='30'>EBR</td><td height='30'> Logica:  " + convertName(ebr.NombreParticion[:]) + " " + strconv.Itoa(int(PorcentajeUtilizao)) + "%</td>")
						cont++
					}
				}
				if (Particiones[i].TamanioTotal - EspacioUtilizado) > 0 {
					PorcentajeUtilizao = (float64(TamanioDisco-EspacioUtilizado) / float64(TamanioDisco)) * 100
					buffer.WriteString("<td height='30' width='100%'>Libre: " + strconv.Itoa(int(PorcentajeUtilizao)) + "%</td>")
				}
			}
			buffer.WriteString("</tr>\n")
			buffer.WriteString("</table>\n</td>")
		}
	}
	if (TamanioDisco - EspacioUtilizado) > 0 {
		PorcentajeUtilizao = (float64(TamanioDisco-EspacioUtilizado) / float64(TamanioDisco)) * 100
		buffer.WriteString("<td height='30' width='75.0'>Libre: " + strconv.Itoa(int(PorcentajeUtilizao)) + "%</td>")
	}
	buffer.WriteString("     </tr>\n</table>\n>];\n}")
	var datos string
	datos = string(buffer.String())
	CreateArchivo(path, datos)
	return false
}

func RecorrerListaDisco(id string, ListaDiscos *list.List) (string, string, string) {
	Id := strings.ReplaceAll(id, "14", "")
	fmt.Println("id de recorrer ", Id)
	//NoParticion := Id[1:]
	IdDisco := Id[:1]
	fmt.Println("id de IdDisco ", Id)
	pathDisco := ""
	nombreParticion := ""
	nombreDisco := ""
	for element := ListaDiscos.Front(); element != nil; element = element.Next() {
		var disco DISCO
		disco = element.Value.(DISCO)
		if BytesToString(disco.Id) == IdDisco {
			for i := 0; i < len(disco.Particiones); i++ {
				var mountTemp = disco.Particiones[i]
				//fmt.Println("ID DISCOS MOUNT", mountTemp.Id)
				if mountTemp.Id == id {
					copy(mountTemp.EstadoMKS[:], "1")
					nombreParticion = mountTemp.NombreParticion
					pathDisco = disco.Path
					nombreDisco = disco.NombreDisco
					return pathDisco, nombreParticion, nombreDisco
					break
				}
			}

		}
		element.Value = disco
	}
	return "", "", ""
}

func CreateArchivo(path string, data string) {
	propiedades := strings.Split(path, "/")
	nombreArchivo := propiedades[len(propiedades)-1]
	f, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(data)

	if err2 != nil {
		log.Fatal(err2)
	}
	executeComand("dot -Tpdf " + path + " -o " + nombreArchivo[0:len(nombreArchivo)-4] + ".pdf")
	//executeComand("xdg-open " + nombreArchivo[0:len(nombreArchivo)-4] + ".pdf")
	//executeComand("xdg-open " + path)

}

func convertName(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
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
